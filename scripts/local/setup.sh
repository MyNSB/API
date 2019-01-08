#!/usr/bin/env bash


#!/usr/bin/env bash

# Configs
function close {
    echo "Ensure that you complete the pg_cron installation, further instructions can be found at: \n https://github.com/citusdata/pg_cron#setting-up-pg_cron"
    echo "Optionally you may wish to add a password to the database"
}
trap close EXIT

# colours ~~ :)
CYAN='\033[0;36m'
RED='\033[0;38m'
NC='\033[0m'
userdir=$(eval echo ~$USER)
api_dir = find ${userdir} "mynsb-api"

# -------- Downloads -----------
# PostgreSQL
service postgres status
if [ "$?" -gt "0" ]; then
    echo -e "${CYAN}Installing PostgreSQL...${NC}"
    # Repository setup
    sudo apt-get install wget ca-certificates
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
    sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
    # installation
    sudo apt-get update
    sudo apt-get install postgresql postgresql-contrib
    echo -e "${CYAN}Finished Installing Postgres${NC}"
fi
# Time to install the postgres packages
# Determine if they are there
# Pretty hacky for checking if an extension exists
dump = sudo -u postgres -H -- psql -c "select count(*) from pg_available_extensions where name='pg_cron';"
if echo "${dump}" | grep '[0]' >/dev/null; then
    # pg_cron doesn't exist
    echo -e "${CYAN}Installing pg_cron...${NC}"
    sudo add-apt-repository ppa:libreoffice/ppa
    sudo apt-get -y install postgresql-10-cron
    echo -e "${CYAN}Finished Installing pg_cron...${NC}"
fi

# Golang
if ! (command -v go |  grep -q ^ >/dev/null); then
    echo -e "${CYAN}Installing Golang...${NC}"
    curl -O https://storage.googleapis.com/golang/go1.11.2.linux-amd64.tar.gz
    tar -xvf go1.11.2.linux-amd64.tar.gz
    sudo mv go ${userdir}
    mkdir ${userdir}/workspace
    echo "export GOROOT=\$HOME/go" >> ${userdir}/.bashrc
    echo "export GOPATH=\$HOME/workspace" >> ${userdir}/.bashrc
    echo "export PATH=\$PATH:\$GOROOT/bin:\$GOPATH/bin" >> ${userdir}/.bashrc
    source ${userdir}/.bashrc
    echo -e "${CYAN}Finished Installing Golang...${NC}"
fi





# ----------- Setting Up -----------

# PostgreSQL
mv ${api_dir} ${GOPATH}/src
# install db
sudo -u postgres -H -- psql -c "\i ${GOROOT}/src/mynsb-api/database/setup.sql" >/dev/null
# API
# install the api
go install mynsb-api


echo -e "${CYAN}CONGRATULATIONS! FINALLY INSTALLED MYNSB....${NC}"