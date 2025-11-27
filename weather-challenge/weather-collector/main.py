import os
import time
import json
import requests
import pika
import schedule
from datetime import datetime

# Configurações
RABBITMQ_URI = os.getenv('RABBITMQ_URI', 'amqp://user:password@rabbitmq:5672')
LATITUDE = os.getenv('LATITUDE', '-23.55')
LONGITUDE = os.getenv('LONGITUDE', '-46.63')
QUEUE_NAME = 'weather_queue'

def get_weather_data():
    print("🌤️ Buscando dados climáticos...")
    try:
        url = f"https://api.open-meteo.com/v1/forecast?latitude={LATITUDE}&longitude={LONGITUDE}&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code"
        response = requests.get(url, timeout=10)
        data = response.json()

        current = data.get('current', {})
        condition_map = {0: 'Céu Limpo', 1: 'Parcialmente Nublado', 2: 'Nublado', 3: 'Encoberto', 61: 'Chuva Leve', 63: 'Chuva', 80: 'Chuva Forte'}

        payload = {
            "city": f"Lat: {LATITUDE}, Lon: {LONGITUDE}",
            "temperature": current.get('temperature_2m'),
            "humidity": current.get('relative_humidity_2m'),
            "windSpeed": current.get('wind_speed_10m'),
            "condition": condition_map.get(current.get('weather_code', 0), 'Desconhecido'),
            "capturedAt": datetime.now().isoformat()
        }
        return payload
    except Exception as e:
        print(f"❌ Erro: {e}")
        return None

def publish_to_queue(payload):
    if not payload: return
    try:
        params = pika.URLParameters(RABBITMQ_URI)
        connection = pika.BlockingConnection(params)
        channel = connection.channel()
        channel.queue_declare(queue=QUEUE_NAME, durable=True)
        channel.basic_publish(exchange='', routing_key=QUEUE_NAME, body=json.dumps(payload))
        print(f"📤 Enviado: {payload['temperature']}°C")
        connection.close()
    except Exception as e:
        print(f"❌ Erro RabbitMQ: {e}")

def job():
    data = get_weather_data()
    publish_to_queue(data)

print("🚀 Weather Collector Iniciado!")
job() # Executa a primeira vez
schedule.every(1).minutes.do(job) # Agenda a cada 1 min

while True:
    schedule.run_pending()
    time.sleep(1)
    