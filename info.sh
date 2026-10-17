#!/bin/bash
# info.sh - сведения о системе
echo "Пользователь: $(whoami)"
echo "Дата и время: $(date '+%d.%m.%Y %H:%M:%S')"
echo "Имя хоста:    $(hostname)"
echo "Свободно на диске (/):"
df -h / | awk 'NR==1 || NR==2 {print "  " $0}'
