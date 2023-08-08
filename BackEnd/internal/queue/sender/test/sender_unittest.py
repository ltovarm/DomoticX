#! /usr/bin/env python3

import unittest
import pika
import json
import os
import hardwareSmoke as smoke
import time

class TestRabbitMQSender(unittest.TestCase):
    def setUp(self):

        # Get the absolute path of the test script directory
        current_directory = os.path.dirname(os.path.abspath(__file__))
        config_file_path = os.path.join(current_directory, 'config.json')

        # Load configuration from the JSON file
        with open(config_file_path, 'r') as config_file:
            config = json.load(config_file)

        # Set up the connection to RabbitMQ (make sure RabbitMQ is running)
        # Configurar los parámetros de conexión
        self.connection = pika.BlockingConnection(pika.URLParameters(config['queue_url']))
        self.channel = self.connection.channel()

        # Define the test queue
        self.queue_name = config['queue_name']
        self.channel.queue_declare(queue=self.queue_name, durable=True)

        smoke.test_sendTCP()
        time.sleep(2)

    def tearDown(self):
        # Close the connection at the end of the tests
        self.connection.close()

    def test_sender(self):

        # Get the absolute path of the test script directory
        current_directory = os.path.dirname(os.path.abspath(__file__))
        data_file_path = os.path.join(current_directory, 'payload_tcp.json')

        # Consume the message from the queue
        method_frame, header_frame, body = self.channel.basic_get(queue=self.queue_name, auto_ack=True)

        # Ensure the message is not null
        self.assertIsNotNone(body, "The message was not received correctly from the queue.")

        # Decode the JSON message
        received_message = json.loads(body)

        # Load configuration from the JSON file
        with open(data_file_path, 'r') as data_file:
            message = json.load(data_file)

            # Perform assertions with the received message
            self.assertEqual(message['id'], received_message['id'], "The content of the message is not as expected.")
            self.assertEqual(message['ndata'], received_message['ndata'], "The content of the message is not as expected.")
            for i in range(len(message['datajson'])):
                self.assertEqual(message['datajson'][i]['id'], received_message['datajson'][i]['id'], "The content of the message is not as expected.")
                self.assertEqual(message['datajson'][i]['msgtype'], received_message['datajson'][i]['msgtype'], "The content of the message is not as expected.")
                self.assertEqual(message['datajson'][i]['data'], received_message['datajson'][i]['data'], "The content of the message is not as expected.")
            
if __name__ == '__main__':
    unittest.main(verbosity=2)
