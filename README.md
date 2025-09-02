<h1 align="center">WB L0.1<p>WB Order Service</p></h1>

<p align="center">
  <img src="https://media4.giphy.com/media/v1.Y2lkPTc5MGI3NjExN3ExNG00M2didG5iMW12azJqM2p4cDEweXgwN2Fic25pbDAxZmM4aSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/lIzAEoZEn571u/giphy.gif" alt="alt" width="300" />
</p>

## 🛠️ Использованный стек

- **Golang 1.23** — основной сервис (сервер API, consumer, producer)  
- **PostgreSQL** — база данных, хранилище заказов  
- **Kafka (Confluent)** — брокер сообщений (3 ноды + Zookeeper)  
- **Kafka UI** — удобная обёртка для дебага топиков  
- **Docker + docker-compose** — контейнеризация  
- **Makefile** — быстрые команды (`make build`, `make launch`, `make produce`)  
- **HTML/CSS/JS (embed)** — лёгкий веб-интерфейс для просмотра заказов 

## 🚀 Инструкция по запуску

1. Склонируй репозиторий  
   ```bash
   git clone <repo_url>
   cd <folder_name>
   ```

2. Запусти Docker

3. Собери и подними окружение:  
   ```bash
   make build
   make launch
   ```

4. Проверь логи (Kafka, Postgres, consumer, server):  
   ```bash
   make logs
   ```

5. Отправь тестовые данные в брокер:  
   ```bash
   make produce
   ```
   ✅ После этого consumer загрузит их в Postgres.

6. Открой в браузере:  
   - UI сервиса: [http://localhost:8080](http://localhost:8080)  
   - Kafka UI: [http://localhost:8081](http://localhost:8081)  

7. Попробуй найти заказ:  
   - Вводи `order_uid` из `sample.json` (например: `test1000`) в форму UI  
   - Или напрямую через API:  
     ```bash
     curl http://localhost:8080/orders/test1000 | jq
     ```


## 📝 Примечания

- Кэш реализован примитивно через map в Golang. Он пополняется после обращения по `order_uid` и при перезапуске сервиса последними 10 добавленными заказами. 
- Если Kafka при старте выдаёт `ApiVersion … Disconnected` — это нормально: брокеры могут быть не сразу готовы, клиент переподключится. 
