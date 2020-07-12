# A Simple full-stack project utilizing Python, Go and Vue.js

This is a simple full-stack project using python to create an populate both sql and Nosql databases with multiple tables. Then used Go to create an api to serve these tables and finally a Vue.js simple web interface to display and filter the data.

## Before trying to run this project make sure you have installed:
- Python
- Go
- Vue.js
- MongoDB 
- Postgres DB

By default the used port for the frontend is 8080 and for the api is 8000, If you changed the api's port make sure to change it in the frontend as will be explained below.

## Steps to run this project:
1) Make sure that both MongoDB and Postgres DB are running.
2) Navigate into the database_migration folder and install the python dependencies using 
```
pip3 install -r requirments.txt
```
3) Run the Python script using
```
python3 database_migration.py
```
then follow the prompts.

4) Navigate to the Go restapi files located at restapi/src/github.com/karimkhalifa/restapi/ and install dependencies using
```
go get ./...
```
5) Run the api server using
```
go run .
```
6) Navigate to the Vue.js front end at frontend/ and install dependencies using
```
npm install
```
7) The front end uses by default port 8080 and the api uses port 8000, if you used another port for the api you will need to change it in the orders.vue file under frontend/src/views/, you will need to change the baseURI and the inUseURI to match your selected api port

8) Finally, run
```
npm run serve
```
## You can:
1) search for items by their names or by order name. 
2) you can filter by dates from the date picker.
3) you can do both at the same time.

Pagination, searching and filtering is done server-side.
