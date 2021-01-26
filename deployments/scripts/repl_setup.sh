#!/bin/bash

docker-ip() {
    docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$@"
}

MS_STATUS=`docker exec db-master sh -c 'export MYSQL_PWD=root_password; mysql -u root -e "SHOW MASTER STATUS"'`
CURRENT_LOG=`echo $MS_STATUS | awk '{print $6}'`
CURRENT_POS=`echo $MS_STATUS | awk '{print $7}'`

echo $MS_STATUS
start_slave_stmt="CHANGE MASTER TO MASTER_HOST='db-master',MASTER_USER='otus_social_slave',MASTER_PASSWORD='otus_social_passwd',MASTER_LOG_FILE='$CURRENT_LOG',MASTER_LOG_POS=$CURRENT_POS; START SLAVE;"
#start_slave_stmt="CHANGE MASTER TO MASTER_HOST='db-master',MASTER_USER='otus_social_slave',MASTER_PASSWORD='otus_social_passwd',MASTER_AUTO_POSITION=1; START SLAVE;"
start_slave_cmd='export MYSQL_PWD=root_password; mysql -u root -e "'
start_slave_cmd+="$start_slave_stmt"
start_slave_cmd+='"'
echo $start_slave_cmd
docker exec db-slave1 sh -c "$start_slave_cmd"
docker exec db-slave2 sh -c "$start_slave_cmd"

docker exec db-slave1 sh -c "export MYSQL_PWD=root_password; mysql -u root -e 'SHOW SLAVE STATUS \G'"
docker exec db-slave2 sh -c "export MYSQL_PWD=root_password; mysql -u root -e 'SHOW SLAVE STATUS \G'"