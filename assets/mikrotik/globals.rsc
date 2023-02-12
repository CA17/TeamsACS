#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Get Serial Number.
# routerboard or chr
# Usage: :global sn [ $GetSN ]
:global GetSN;
:if (!any $GetSN) do={:global GetSN do={
    :local sn1 [/system routerboard get serial-number];
    # sn for CHR
    :if ([:len $sn1] = 0) do={
        :set sn1 [/system license get system-id];
    }

    :return $sn1;
}}

# Either first (if not nil) Or second
# Usage: $Either $val $defaultval
:global EitherOr;
:if (!any $EitherOr) do={:global EitherOr do={
    :if ([ :len $1 ] > 0) do={ :return $1; } else={ :return $2; }
}}

# Character Replace
# to replace all the occurrences of a certain segment within a given string
# Usage: $CharacterReplace $str $segment $replacement
:global CharacterReplace;
:if (!any $CharacterReplace) do={:global CharacterReplace do={
    :local String [ :tostr $1 ];
    :local ReplaceFrom [ :tostr $2 ];
    :local ReplaceWith [ :tostr $3 ];
    :local Return "";

    :if ($ReplaceFrom = "") do={
        :return $String;
    }

    :while ([ :typeof [ :find $String $ReplaceFrom ] ] != "nil") do={
        :local Pos [ :find $String $ReplaceFrom ];
        :set Return ($Return . [ :pick $String 0 $Pos ] . $ReplaceWith);
        :set String [ :pick $String ($Pos + [ :len $ReplaceFrom ]) [ :len $String ] ];
    }

    :return ($Return . $String);
}}

# Split
:global Split;
:if (!any $Split) do={:global Split do={
    :local String [ :tostr $1 ];
    :local Splitter [ :tostr $2 ];
    :local Return;

    :global CharacterReplace;

    :if ($Splitter = "") do={
        :return $String;
    }

    :set String [ $CharacterReplace $String $Splitter "," ];
    :set Return [ :toarray $String ];

    :return $Return;
}}

# Left
:global Left;
:if (!any $Left) do={:global Left do={
    :global Split;

    :local tmpArr [ $Split $1 $2 ];
    :return ($tmpArr->0);
}}

# Trim
:global Trim;
:if (!any $Trim) do={:global Trim do={
    :local Blank [ :tostr $2 ];
    :global CharacterReplace;

    :if ($Blank = "") do={
        :set Blank " ";
    }

    :return [ $CharacterReplace $1 $Blank "" ];
}}

# Get Timestamp
# Return: yyyymmddhhmiss
:global GetTimestamp;
:if (!any $GetTimestamp) do={:global GetTimestamp do={
    :global CharacterReplace;

    :local tm [ /system clock get time ];
    :local dt [ /system clock get date ];
    :local monMap {"jan"="01";"feb"="02";"mar"="03";"apr"="04";"may"="05";"jun"="06"; \
                    "jul"="07";"aug"="08";"sep"="09";"oct"="10";"nov"="11";"dec"="12"};

    :local year [:pick $dt 7 11];
    :local mon [:pick $dt 0 3];
    :local day [:pick $dt 4 6];
    :set dt ($year . ($monMap->$mon) . $day);
    # :set dt [ $CharacterReplace $dt $mon ($monMap->$mon) ];

    :set tm [ $CharacterReplace $tm ":" "" ];

    :return ($dt . $tm);
}}

# ToJson
# Turns array to JSON format.
# Usage: $ToJson $arr
:global ToJson;
:if (!any $ToJson) do={:global ToJson do={
    :local arr $1;
    :local json "";

    :if ([ :typeof $arr] != "array") do={
        return "";
    }

    :foreach k,v in=$arr do={
        :if ($json = "") do={
            :set json "{";
        } else={
            :set json "$json,";
        }
        :set json ("$json\"$k\":\"$v\"");
    }

    :set json "$json}";
    :return $json;
}}

# GetIPOnIface
# Get IP address on the given interface in CIDR form.
# Usage: $GetIPOnIface iface
:global GetIPOnIface;
:set GetIPOnIface do={
    :local iface [ :tostr $1 ];
    :local ifip "";

    :if ([ :len $iface ] = 0) do={
        :return "";
    }

    # loop every 5s for about 1m to wait for the address
    # the reason for 1m is that netwatch should be up in less than 1m
    :local I 0;
    :while ($I <= 12) do={
        :do {
            /ip address {
                :set ifip [ get [ find where interface ~$iface ] address ];
                :return $ifip;
            }
        } on-error={ :log info "$iface IP is not available yet. $I tries."};

        :set I ($I + 1);
        :delay 5s;
    }

    :return $ifip;
}

# GetAssetID
# Get asset ID, in the form of "XXX-XXX", default to be the last two value of the Managelink IP
# Usage: $GetAssetID manip
:global GetAssetID;
:set GetAssetID do={
    :local manip [ :tostr $1 ];
    :local maniparr;

    :global Split;
    :global EitherOr;

    :set manip [ $EitherOr $manip "0.0.0.0/0" ];
    :set maniparr [ $Split $manip "/" ];
    :set manip ($maniparr->0);

    :set maniparr [ $Split $manip "." ];
    :local id1 ($maniparr->2);
    :local id2 ($maniparr->3);

    :set id1  [ :pick [ :tostr (1000 + $id1) ] 1 4];
    :set id2  [ :pick [ :tostr (1000 + $id2) ] 1 4];

    :return ("$id1-$id2");
}

# CalcNet
# Calculate Network for specific ip address
# Usage: $CalcNet $ipaddr
:global CalcNet;
:set CalcNet do={
    :local ipaddr [ :tostr $1 ];
    :local Return [ :toarray "" ];

    :if ($ipaddr = "") do={
        :error "Parameter ipaddr is required.";
    }

    :local slash [:find $ipaddr "/"];
    :local tmpip [:pick $ipaddr 0 $slash];
    :local tmpmask [:pick $ipaddr ($slash+1) [:len $ipaddr]];
    :local CIDRmask (255.255.255.255 << (32-$tmpmask) );
    :local tmpnet ($tmpip & $CIDRmask);

    :set ($Return->"ip") $tmpip;
    :set ($Return->"mask") $tmpmask;
    :set ($Return->"CIDRmask") $CIDRmask;
    :set ($Return->"net") $tmpnet;

    :return $Return;
}
