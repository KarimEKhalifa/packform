import os
import glob2
import json
import pymongo
import psycopg2
import pandas as pd
from sqlalchemy import create_engine
from sqlalchemy.pool import NullPool


def CreateDB(db_name, pgUser, mongoURI):
    try:
        postgres_con = psycopg2.connect(database="postgres",user=pgUser)
        postgres_con.autocommit = True
        cursor = postgres_con.cursor()
        res = cursor.execute(f'''CREATE database {db_name}''')
        postgres_con.close()
    except:
        print("Postgres database exists") 
    mongo_client = pymongo.MongoClient(f"mongodb://{mongoURI}/")
    engine = create_engine(f"postgres:///{db_name}",poolclass=NullPool)
    return (db_name,engine,mongo_client)

def CreateTables(db_name,engine,mongo_client,files):
    for f in files:
        table_name = f.rsplit("-",1)[1].replace(".csv","").strip()
        df = pd.read_csv(f)
        if "Mongo" in f:
            InsertMongo(db_name,mongo_client,df,table_name)
        elif "Postgres" in f:
            try:
                df.to_sql(table_name, engine, index=False)
            except:
                print(f"Table {table_name} already exists")

def InsertMongo(db_name,myclient,df,table_name):
    mydb = myclient[db_name]
    collection = mydb[table_name]
    records = json.loads(df.T.to_json()).values()
    for record in records:
        existing_record = collection.find_one(record)
        if not existing_record:
            collection.insert_one(record)
        else:
            collection.replace_one(existing_record,record,upsert=True)

if __name__ == "__main__":
    print("Before running this script make sure that:")
    print("1- MongoDB is running.")
    print("2- Postgres is running.")
    print("3- Filenames must include Postgres or Mongo in their names seperated by '-' to identify which db to create for the respective table." )
    print("4- Folder including the csv files should be in the same directory as the script" )
    pgUser = input("Enter postgres username: ")
    mongoURI = input("Enter Mongodb URI ex. localhost:27017 : ")
    db_name = input("Enter database name to create and populate: ")
    csv_folder = input("Enter folder name: ")
    
    try:
        csvs = glob2.glob(os.path.join(os.getcwd(),f"{csv_folder}/*.csv"))
        CreateTables(*CreateDB(db_name,pgUser,mongoURI),csvs)
        print("Finished successfully...")
    except:
        print("An error has occured!")





