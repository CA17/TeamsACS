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

# local vars for future use
:local metrics;
:local httpdata "";

# obtain main infos
:local Res [ /system resource get ];

# Create metrics array
:set metrics {
    "system.info.sn"=[ $sn ];
    "system.info.timestamp"=[ $timestamp ];
    "system.cpu.load"=($Res->"cpu-load");
    "system.memory.total"=($Res->"total-memory");
    "system.memory.free"=($Res->"free-memory");
    "system.disk.total"=($Res->"total-hdd-space");
    "system.disk.free"=($Res->"free-hdd-space");
    "system.firewall.connections"=[ /ip firewall connection print count-only ];
}

# Datadog JSON data to parse with Datadog API post-timeseries-points
:set httpdata [ $ToJson $metrics ];
:log info "device_metric: $httpdata";

# Call API via POST
:local ret [ /tool fetch mode=https http-method=post http-header-field="Content-Type:application/json,Authorization:Bearer $TeamsacsApiToken" \
    http-data=$httpdata url="$TeamsacsApiServer/api/v1/device/metric" keep-result=no as-value ];

:log info ("device_metric: " . [ :tostr $ret ]);


