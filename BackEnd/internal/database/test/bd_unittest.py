#! /usr/bin/env python3

import unittest
import psycopg2
import os
import json
import sqlparse

class TestDatabaseCreation(unittest.TestCase):
    def setUp(self):

        # Get the absolute path of the test script directory
        current_directory = os.path.dirname(os.path.abspath(__file__))
        config_file_path = os.path.join(current_directory, 'config.json')
        
        # Load configuration from the JSON file
        with open(config_file_path, 'r') as config_file:
            config = json.load(config_file)

        # Connect to the PostgreSQL database
        self.connection = psycopg2.connect(
            host=config["host"],
            port=config["port"],
            user=config["user"],
            password=config["password"],
            database=config["database"]
        )

        config_db_path = os.path.join(current_directory, '../init_db.sql')

        # Load and execute the init_db.sql file to create the database
        with open(config_db_path, 'r') as sql_file:
            sql_commands = sql_file.read()

        # Separate SQL statements in a list
        commands = sqlparse.split(sql_commands)

        # Execute each SQL statement separately
        with self.connection.cursor() as cursor:
            for command in commands:
                # If the statement is not empty and does not begin with "\", execute it
                if command.strip() and not command.strip().startswith('\\'):
                    cursor.execute(command)
            self.connection.commit()

    def tearDown(self):
        self.connection.close()

    def test_database_creation(self):

        # Check if the database exists in PostgreSQL
        with self.connection.cursor() as cursor:
            cursor.execute("SELECT datname FROM pg_database WHERE datname='house';")
            result = cursor.fetchone()

        # Perform the assertion
        self.assertIsNotNone(result, "The database has not been created correctly in PostgreSQL.")

    def test_table_creation(self):
        
        # Verify that table 'temperatures' exists in the schema
        with self.connection.cursor() as cursor:
            cursor.execute("""
                SELECT column_name, data_type 
                FROM information_schema.columns 
                WHERE table_name = 'temperatures';
            """)
            columns = cursor.fetchall()
        
        # Verify that the table has the 'id' and 'data' columns with the appropriate types.
        expected_columns = [('id', 'integer'), ('data', 'jsonb')]
        for column_name, data_type in expected_columns:
            self.assertIn((column_name, data_type), columns)


if __name__ == '__main__':
    unittest.main()
