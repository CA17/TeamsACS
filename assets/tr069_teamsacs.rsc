# Install certificate
:global acsCaCertTxt "-----BEGIN CERTIFICATE-----
MIIFyzCCA7OgAwIBAgIUahBkSL5OKA5T50NBjgRBYbugGLEwDQYJKoZIhvcNAQEL
BQAwdTELMAkGA1UEBhMCQ04xETAPBgNVBAgMCFNoYW5naGFpMRQwEgYDVQQKDAt0
b3VnaHJhZGl1czEWMBQGA1UEAwwNVG91Z2hyYWRpdXNDQTElMCMGCSqGSIb3DQEJ
ARYWbWFzdGVyQHRvdWdocmFkaXVzLm5ldDAeFw0yMzAyMDYwODA0MzlaFw0zMzAy
MDMwODA0MzlaMHUxCzAJBgNVBAYTAkNOMREwDwYDVQQIDAhTaGFuZ2hhaTEUMBIG
A1UECgwLdG91Z2hyYWRpdXMxFjAUBgNVBAMMDVRvdWdocmFkaXVzQ0ExJTAjBgkq
hkiG9w0BCQEWFm1hc3RlckB0b3VnaHJhZGl1cy5uZXQwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQDaZIRZpHp45SyI7lAcMOxEKUCJl4H/8OX+zZG2jbGp
f+uaDvWxNeUJIEongT1xYAzDjP3w+SaE4jiTKjHFo+099CgGkUlPfCvGfVMJabnH
jOdeddJmgaRHQWrGlG2eywF296dxXN31fG6yZC5i60tZn9exrVRVMT49Y6mNow6b
KpNwYww05PnrfzbSRn3YFIrpho4qcek4KexfFAJHg9dm0oBSBy/rGNqls7TxnnXs
e94IUN6w7oF3MZrxATqVDGPVvO/VvLIB3Vpm23LxqMOmEH0P4DTxaqraBYJlwSBu
8PYwxFvVx66gQVMyFISPJs+1/zO/QM2ZmJpp+ecAr/6D8qK182tJmpLCjODw7tiA
T+UuRGoGhuwdY4uzmEl+Drd469JLjScAPNr1jAGl6mDeS9O4+L+QnyAIyzeSu7Ty
JwlzPTfDz0nCR+S/zgfsrtHRRzwnVhTKBLZP5YdMf6f5zSPBaqJnSpOrDy1/GLsK
c9iobdFR49V6rby5gwhUvK4wnlIV3T4ZaU+ces3/lSFBlKAZh28fLZgNczB/cTG0
ADdJtnfFs6/xL2vXT6OEHThDW3BSmEf57DC7B0z/pHWKuwv+alUHmUjLP6wCFi5W
YvwSqnxqomM7dsHFIjMu6cXu2gz146ASml6TVCVDk3imyENRsODpnmlU39L2tG+W
kQIDAQABo1MwUTAdBgNVHQ4EFgQUYbwxKZj8sX5UqJOZNUX3JXrMsu0wHwYDVR0j
BBgwFoAUYbwxKZj8sX5UqJOZNUX3JXrMsu0wDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAgEAv/aLkoiLnzsMzWCJVh5WPxDkLyZc1tmg0l040kUB609e
pWvWfU1rQ0mAQ/V74E7sfPIoT/HLtqtrcqWwLiYyTTySgkMLWJex/qCO+kUgvz5+
Oto8/QgGc8pT7u10OGHGflYzZycWsRB4TkkjJW0ejjKjOxJhmpKxFHnsvP1Qj4Qa
j6M91UIiDiUu8eVFbnvXVWd/oIawRquV/PqeCt1W9paIHqGymkFR1TwavgUc0oCJ
382sgNcsgQqIskyLjajFgidoWMCDxeIltjn70kj5kDAnqq10HCzGctbCSdRCF9iR
20UULJn0ynLL/g6baqN0xu+nDDkyoPI9rYlZsZGNZ2CbS07YoZUqYv2EIx4l4cJk
5A3L+QsmSU8uuAGUD6RkiX4FEtZjhtSY7q+f2kqeu2J7AvUA899IWBD6aRfJ571l
irrQcZ3mVqtOUOyFlSHDAVrT3F+xOcimn49MhgPnROi+SRE2xgcWA4eO36hskvJF
XsSBBZWd8R2rtbLtm9Fqun0ULRYDF4f7g8RUz4f/GdFP63ghs4ZY14c0zwKid2Fa
4BEoKtqKMC/2Ua0UBGBF40GL2kyW3aRsUub/hnSs0jzMCw9GZ3t54Wr2Gy85PSdL
83NezyAjASlkhKZXrRjOz143jTULnW5QkX+bSk4PMnOZ+fqqMYiGKby2+ndCN4A=
-----END CERTIFICATE-----";

/file print file=tmp_acs_ca_cert.txt;
delay 2;
/file set tmp_acs_ca_cert.txt contents=$acsCaCertTxt;
/certificate import file-name=tmp_acs_ca_cert.txt passphrase="";
/file remove tmp_acs_ca_cert.txt;



#Get serial-number
:local sn;
:set sn [/system routerboard get serial-number];
:if ([:len $sn]=0) do={
    :set sn [/system license get system-id];
}

/tr069-client set acs-url="https://acs.teamsacs.cc:1989" enabled=yes \
username="$sn" password="examplesecurepassword" periodic-inform-interval=30s

:local setupdate [/system clock get date];
:local setuptime [/system clock get time];

:local note ("# Device Info \r\
    \n1. Serial Number: $sn \r\
    \n2. TR069 Setup Time: $setupdate $setuptime \r\
    \n");

/system note set note=$note;