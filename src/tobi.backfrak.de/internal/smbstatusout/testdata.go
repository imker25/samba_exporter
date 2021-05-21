package smbstatusout

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

const LockDataOneLine = `
Locked files:
Pid          User(ID)   DenyMode   Access      R/W        Oplock           SharePath   Name   Time
--------------------------------------------------------------------------------------------------
1120         1080       DENY_NONE  0x80        RDONLY     NONE             /usr/share/data   .   Sun May 16 12:07:02 2021`

const LockData0Line = `
Locked files:
Pid          User(ID)   DenyMode   Access      R/W        Oplock           SharePath   Name   Time
--------------------------------------------------------------------------------------------------`

const LockData4Lines = `Locked files:
Pid          User(ID)   DenyMode   Access      R/W        Oplock           SharePath   Name   Time
--------------------------------------------------------------------------------------------------
1120         1080       DENY_NONE  0x80        RDONLY     NONE             /usr/share/data    .   Sun May 16 12:07:02 2021
1120         1080       DENY_NONE  0x80        RDONLY     NONE             /usr/share/foto    .   Mon May 17 06:39:38 2021
1120         1080       DENY_NONE  0x80        RDONLY     NONE             /usr/share/film    .   Mon May 17 07:09:38 2021
1120         1080       DENY_NONE  0x80        RDONLY     NONE             /usr/share/music   .   Fri May 21 18:46:29 2021`

const ShareDataOneLine = `
Service      pid     Machine       Connected at                      Encryption   Signing     
---------------------------------------------------------------------------------------------
IPC$         1119    192.168.1.242  Sun May 16 11:55:36 AM 2021 CEST -            -           `

const ShareData0Line = `
Service      pid     Machine       Connected at                      Encryption   Signing     
---------------------------------------------------------------------------------------------`

const ShareData4Lines = `
Service      pid     Machine       Connected at                      Encryption   Signing     
---------------------------------------------------------------------------------------------
IPC$         1119    192.168.1.242  Sun May 16 11:55:36 AM 2021 CEST -            -           
foto         1121    192.168.1.243  Mon May 17 10:56:56 AM 2021 CEST -            -           
film         1117    192.168.1.244  Tue May 18 09:52:38 AM 2021 CEST -            -           
musik        1117    192.168.1.245  Fri May 21 06:46:29 PM 2021 CEST -            -           `

const ProcessDataOneLine = `
Samba version 4.11.6-Ubuntu
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing              
----------------------------------------------------------------------------------------------------------------------------------------
1117    1080         ssl-cert     192.168.1.242 (ipv4:192.168.1.242:42296)  SMB3_11           -                    partial(AES-128-CMAC)`

const ProcessData4Lines = `
Samba version 4.11.6-Ubuntu
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing              
----------------------------------------------------------------------------------------------------------------------------------------
1117    1080         ssl-cert     192.168.1.242 (ipv4:192.168.1.242:42296)  SMB3_11           -                    partial(AES-128-CMAC)
1119    1080         ssl-cert     192.168.1.243 (ipv4:192.168.1.243:47510)  SMB3_11           -                    partial(AES-128-CMAC)
1120    1080         ssl-cert     192.168.1.244 (ipv4:192.168.1.244:47512)  SMB3_11           -                    partial(AES-128-CMAC)
1121    1080         ssl-cert     192.168.1.245 (ipv4:192.168.1.245:47514)  SMB3_11           -                    partial(AES-128-CMAC)`

const ProcessData0Lines = `
Samba version 4.11.6-Ubuntu
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing              
----------------------------------------------------------------------------------------------------------------------------------------`
