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
1120         1080       DENY_NONE  0x80        RDONLY     NONE             /usr/share/music   .   Sun Oct  1 12:39:21 2022`

const LockDataNoData = `No locked files`

const LockDataNoDataV4_17_7 = `No locked files
`

const LockDataEmpty = `  
  
`

const LockDataCluster = `Locked files:
Pid          Uid        DenyMode   Access      R/W        Oplock           SharePath   Name   Time
--------------------------------------------------------------------------------------------------
1:55399      1001       DENY_NONE  0x12019f    RDWR       LEASE(RWH)       /lfsmnt/dst01   share/data/data1/Clip/792_2134.MXF 48000_11.pek   Tue Apr  4 14:23:18 2023
1:55399      1001       DENY_WRITE 0x120089    RDONLY     LEASE(RWH)       /lfsmnt/dst01   share/data/data1/Clip/Clip0005.MXF 48000_21.pek   Tue Apr  4 14:26:09 2023
1:55399      1001       DENY_WRITE 0x120089    RDONLY     LEASE(RWH)       /lfsmnt/dst01   share/data/data1/folder/dir/100MEDIA/DJI_0177.MOV   Tue Apr  4 14:32:01 2023
1:55399      1001       DENY_NONE  0x120089    RDONLY     LEASE(RWH)       /lfsmnt/dst01   share/data/data1/folder/dir/100MEDIA/DJI_0177.MOV   Tue Apr  4 14:32:01 2023
1:19801      1001       DENY_NONE  0x100081    RDONLY     NONE             /lfsmnt/dst01   share/dir/data/test_The_Whole.mov  Tue Apr  4 03:17:50 2023
3:57086      1001       DENY_NONE  0x120089    RDONLY     LEASE(RWH)       /lfsmnt/dst01   share/data2/CLIPS001/CC0639/CC063904.MXF   Tue Apr  4 08:11:39 2023
1:55399      1001       DENY_WRITE 0x120089    RDONLY     LEASE(RWH)       /lfsmnt/dst01   share/test.wav 48000.pek   Tue Apr  4 14:13:28 2023`

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
musik        1117    192.168.1.245  Fri Nov 5 11:07:13 PM 2021 CET   -            -           `

const ShareData4LinesWithSpacesInName = `
Service      pid     Machine       Connected at                      Encryption   Signing     
---------------------------------------------------------------------------------------------
test share        4642    127.0.0.1     Mon May 31 17:23:44 2021 UTC     -            -           
IPC$ admin share  4642    127.0.0.1     Wed Jun  2 21:32:31 2021 UTC     -            -    
a b c d e f g h i 4642    127.0.0.1     Wed Jun  2 21:32:31 2021 UTC     -            -    
musik             1117    192.168.1.245  Mo Sep 19 18:34:17 2022 CEST    -            -        `

const ShareDataDifferentTimeStampLines = `
Service      pid     Machine       Connected at                     Encryption   Signing     
---------------------------------------------------------------------------------------------
test         4642    127.0.0.1     Mon May 31 17:23:44 2021 UTC     -            -           
IPC$         4642    127.0.0.1     Wed Jun  2 21:32:31 2021 UTC     -            -    
musik        1117    192.168.1.245  Mo Sep 19 18:34:17 2022 CEST    -            -        `

const ShareDataCluster = `Samba version 4.9.5-Debian
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing
----------------------------------------------------------------------------------------------------------------------------------------
1:19801 nobody       nogroup      10.63.0.36 (ipv4:10.63.0.36:53407)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:19801 nobody       nogroup      10.63.0.36 (ipv4:10.63.0.36:53407)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:25648 nobody       nogroup      10.63.0.81 (ipv4:10.63.0.81:49591)        SMB3_11           -                    -
1:25648 nobody       nogroup      10.63.0.81 (ipv4:10.63.0.81:49591)        SMB3_11           -                    -
1:25648 nobody       nogroup      10.63.0.81 (ipv4:10.63.0.81:49591)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:19801 nobody       nogroup      10.63.0.36 (ipv4:10.63.0.36:53407)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
2:42597 nobody       nogroup      10.63.1.55 (ipv4:10.63.1.55:57033)        SMB3_11           -                    -
1:19801 nobody       nogroup      10.63.0.36 (ipv4:10.63.0.36:53407)        SMB3_11           -                    -
2:42597 nobody       nogroup      10.63.1.55 (ipv4:10.63.1.55:57033)        SMB3_11           -                    -`

const ShareDataEmpty = `  
  
`

const ProcessDataOneLine = `
Samba version 4.11.6-Ubuntu
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing              
----------------------------------------------------------------------------------------------------------------------------------------
1117    1080         117          192.168.1.242 (ipv4:192.168.1.242:42296)  SMB3_11           -                    partial(AES-128-CMAC)`

const ProcessData4Lines = `
Samba version 4.11.6-Ubuntu
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing              
----------------------------------------------------------------------------------------------------------------------------------------
1117    1080         117          192.168.1.242 (ipv4:192.168.1.242:42296)  SMB3_11           -                    partial(AES-128-CMAC)
1119    1080         117          192.168.1.243 (ipv4:192.168.1.243:47510)  SMB3_11           -                    partial(AES-128-CMAC)
1120    1080         117          192.168.1.244 (ipv4:192.168.1.244:47512)  SMB3_11           -                    partial(AES-128-CMAC)
1121    1080         117          192.168.1.245 (ipv4:192.168.1.245:47514)  SMB3_11           -                    partial(AES-128-CMAC)`

const ProcessData0Lines = `
Samba version 4.11.6-Ubuntu
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing              
----------------------------------------------------------------------------------------------------------------------------------------`

const ProcessDataCluster = `Samba version 4.9.5-Debian
PID     Username     Group        Machine                                   Protocol Version  Encryption           Signing
----------------------------------------------------------------------------------------------------------------------------------------
3:57086 nobody       nogroup      10.63.0.41 (ipv4:10.63.0.41:62834)        SMB3_11           -                    -
3:24179 nobody       nogroup      10.63.0.28 (ipv4:10.63.0.28:58968)        SMB3_11           -                    -
1:19801 nobody       nogroup      10.63.0.36 (ipv4:10.63.0.36:53407)        SMB3_11           -                    -
3:24179 nobody       nogroup      10.63.0.28 (ipv4:10.63.0.28:58968)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -
1:55399 nobody       nogroup      10.63.0.11 (ipv4:10.63.0.11:50370)        SMB3_11           -                    -`

const ProcessDataEmpty = `  `
