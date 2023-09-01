#! /usr/bin/env python3

import unittest
import pika
import json
import os
import time
import psycopg2
import websocket
import asyncio
import websockets


class TestConsumer(unittest.TestCase):
    def setUp(self):

        # Get the absolute path of the test script directory
        current_directory = os.path.dirname(os.path.abspath(__file__))
        config_file_path = os.path.join(current_directory, 'config.json')
        data_file_path = os.path.join(current_directory, 'payload_tcp.json')

        # Load configuration from the JSON file
        with open(config_file_path, 'r') as config_file:
            self.config = json.load(config_file)

        # Load configuration from the JSON file
        with open(data_file_path, 'r') as data_file:
            self.message = json.load(data_file)

        # Set up the connection to RabbitMQ (make sure RabbitMQ is running)
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host=self.config['host']))
        self.channel = self.connection.channel()

        # Define the test queue
        self.queue_name = self.config['queue_name']
        self.channel.queue_declare(queue=self.queue_name, durable=True)

        # Connect to database
        self.connection_db = psycopg2.connect(
            host=self.config['host'],
            port=self.config['db_port'],
            database=self.config['db_name'],
            user=self.config['db_user'],
            password=self.config['db_password']
        )

    def tearDown(self):
        # Close the connection at the end of the tests
        self.connection.close()
        self.connection_db.close()

    def test_consumer_send_to_database(self):

        # Publish the msg
        self.channel.basic_publish(exchange='',
                                   routing_key=self.queue_name,
                                   body=json.dumps(self.message))

        query = f"SELECT * FROM {self.config['db_table']} ORDER BY id DESC LIMIT 5;"

        # Execute the query
        with self.connection_db.cursor() as cursor:
            cursor.execute(query)
            rows = cursor.fetchall()

        for i, row in enumerate(reversed(rows)):

            # Decode the JSON message
            received_message = row[1]

            # Perform assertions with the received message
            self.assertEqual(self.message['datajson'][i]['id'], received_message['id'],
                             "The content of the message is not as expected.")
            self.assertEqual(self.message['datajson'][i]['msgtype'],
                             received_message['msgtype'], "The content of the message is not as expected.")
            self.assertEqual(self.message['datajson'][i]['data'], received_message['data'],
                             "The content of the message is not as expected.")

    async def send_data_to_websocket(self):

        # Publish the msg
        self.channel.basic_publish(exchange='',
                                   routing_key=self.queue_name,
                                   body=json.dumps(self.message))

    async def receive_from_websocket(self):
        async with websockets.connect('ws://localhost:8080/ws') as websocket:
            while True:
                received_message = await websocket.recv()
                received_message = json.loads(received_message)
                self.assertEqual(self.message['id'], received_message['id'],
                                 "The content of the message is not as expected.")
                self.assertEqual(self.message['ndata'], received_message['ndata'],
                                 "The content of the message is not as expected.")
                for i in range(len(self.message['datajson'])):
                    self.assertEqual(self.message['datajson'][i]['id'], received_message['datajson']
                                     [i]['id'], "The content of the message is not as expected.")
                    self.assertEqual(self.message['datajson'][i]['msgtype'], received_message['datajson']
                                     [i]['msgtype'], "The content of the message is not as expected.")
                    self.assertEqual(self.message['datajson'][i]['data'], received_message['datajson']
                                     [i]['data'], "The content of the message is not as expected.")

    def test_consumer_send_to_websocket(self):

        # Create a task to receive the message at the WebSocket server
        receive_task = asyncio.ensure_future(self.receive_from_websocket())

        # Create a task to send message to WebSocket server
        send_task = asyncio.ensure_future(self.send_data_to_websocket())


if __name__ == '__main__':
    unittest.main(verbosity=2)
