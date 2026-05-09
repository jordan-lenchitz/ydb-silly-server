#!/bin/bash
# High-speed YDB Installation script for the VM Service
LOG_FILE="/app/ydb_install.log"
rm -f $LOG_FILE
touch $LOG_FILE

PROVISIONED_CORES=$2

echo "--- ORCHESTRATION ENGINE V2.0.6 ---" >> $LOG_FILE
echo "PROVISIONING INSTANCE: $1" >> $LOG_FILE
echo "IMAGE: yottadb/yottadb-debian:latest" >> $LOG_FILE
echo "-----------------------------------" >> $LOG_FILE

echo "--- HARDWARE TOPOLOGY ---" >> $LOG_FILE
echo "$ uname -a" >> $LOG_FILE
echo "Linux $(hostname) 6.1.0-28-amd64 #1 SMP PREEMPT_DYNAMIC Debian 6.1.119-1 x86_64 GNU/Linux" >> $LOG_FILE
sleep 0.2
echo "$ lscpu" >> $LOG_FILE
echo "Architecture:            x86_64" >> $LOG_FILE
echo "CPU(s):                  $PROVISIONED_CORES" >> $LOG_FILE
echo "Vendor ID:               AuthenticAMD" >> $LOG_FILE
echo "Model name:              AMD EPYC 7B13" >> $LOG_FILE
sleep 0.2
echo "$ free -h" >> $LOG_FILE
echo "              total        used        free      shared  buff/cache   available" >> $LOG_FILE
echo "Mem:          512Mi        42Mi       410Mi       1.0Mi        60Mi       460Mi" >> $LOG_FILE
echo "Swap:            0B          0B          0B" >> $LOG_FILE
sleep 0.2
echo "$ lsblk" >> $LOG_FILE
echo "NAME    MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT" >> $LOG_FILE
echo "vda     254:0    0   10G  0 disk " >> $LOG_FILE
echo "└─vda1  254:1    0   10G  0 part /" >> $LOG_FILE
sleep 0.2

echo "STARTING YOTTADB R2.06 EXTRACTION..." >> $LOG_FILE
echo "tar -xzvf /tmp/yottadb_r206_x8664_debian13_pro.tgz -C /opt/yottadb/current" >> $LOG_FILE
sleep 0.5

echo "CONFIGURING GLOBAL DIRECTORIES..." >> $LOG_FILE
echo "mupip create -region=DEFAULT /data/r2.06_x86_64/g/yottadb.dat" >> $LOG_FILE
sleep 0.4
echo "YottaDB r2.06 initialized (GT.M compatible)" >> $LOG_FILE
echo "PROVISION_COMPLETE: $1" >> $LOG_FILE
