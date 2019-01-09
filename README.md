<img src="https://github.com/MyNSB/Android/blob/master/app/src/main/res/mipmap-xxxhdpi/mynsb_logo.png" width="200"/> &nbsp; 
# &nbsp;&nbsp; MyNSB API

[![Travis](https://img.shields.io/travis/MyNSB/API.svg?style=flat-square)](https://travis-ci.org/MyNSB/API)
[![Code Climate](https://img.shields.io/codeclimate/maintainability/MyNSB/API.svg?style=flat-square)](https://codeclimate.com/github/MyNSB/API)
[![Code Climate](https://img.shields.io/codeclimate/issues/github/MyNSB/API.svg?style=flat-square)](https://codeclimate.com/github/MyNSB/API)
[![license](https://img.shields.io/github/license/MyNSB/API.svg?style=flat-square)]()
[![GitHub release](https://img.shields.io/github/release/MyNSB/API.svg?style=flat-square)]()
[![GitHub contributors](https://img.shields.io/github/contributors/MyNSB/API.svg?style=flat-square)](https://github.com/MyNSB/API)


## Setup
The repository includes a simple installation script located [here](scripts/local/setup.sh). 
 - <b>Installation Instructions</b><br>
    ```console
       foo@bar:~$ git clone https://github.com/MyNSB/API.git mynsb-api
       foo@bar:~$ cd mynsb-api/scripts/local
       foo@bar:~/mynsb-api/scripts/local$ sudo sh setup.sh
    ```
    - <b>Configuration</b>
        - If you PostgreSQL installation is running on a port that is not the default: `5432` or if your PostgreSQL installation is not local, you may want to configure the database details file located [here](database/details.txt)    

## Development
All files are located in the `internal` folder.
 - <b>IDE setup</b>
    - The recommended IDE for development is JetBrain's Goland, configuration files for this IDE can be found at: `development/IDE/.idea`. <br>
    - We also recommend that you install the `govendor` tool found: [here](https://github.com/kardianos/govendor) in order to properly manage project dependencies <br>
 - <b>Testing</b>
    - Due to the nature of the application it is recommended that you use a simple API client such as Postman in order to test your code.
 - <b>Contributing</b>
    - Fork and make a pull request :\)          
    
## Usage
 - If the API has been added to your `GOPATH` and `GOBIN` has been added to your `PATH` variable then execution of the API is as follows:
 - ```console
     foo@bar:~$ go install mynsb-api
     foo@bar:~$ mynsb-api
     ```
- This will start a local testing server on port `8080`
- If you are unable to install the API this way then a simple compilation of the source code will work too.
- ```console
    foo@bar:~$ go build $GOPATH/src/mynsb-api/main.go
  ```
  
## Remote Usage
- A remote version of the API may be found at: [https://mynsb.visions.com/api/v1](https://mynsb.visions.com/api/v1)  
- Documentation regarding the general usage of the API can be found at the WIKI section of this repository 
