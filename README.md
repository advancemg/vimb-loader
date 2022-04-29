![GitHub top language](https://img.shields.io/github/languages/top/advancemg/vimb-loader) ![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/advancemg/vimb-loader) ![GitHub all releases](https://img.shields.io/github/downloads/advancemg/vimb-loader/releases/vimb-loader-linux-amd64.zip) ![GitHub](https://img.shields.io/github/license/advancemg/vimb-loader)

## Основные возможности:

* Сервис позволяет работать с ВИМБ через REST API в формате JSON минуя SOAP с нестандартным XML.
* Автоматическая загрузка данных по расписанию. При автоматической загрузке все данные скачиваются на основании бюджетов, бюджеты берутся от текущего дня до конца года.
* Агрегирование данных не предоставляемых ВИМБ.
* Скачанные данные хранятся на S3 и Badger.
* API позволяет вытягивать сохраненные данные формирую динамический запрос.
* В проекте используются RabbitMQ и S3, при указании в config.json хоста 127.0.0.1 или localhost для RabbitMQ и S3
  автоматически запустятся всроенные RabbitMQ и MinIO.
* Для работы всего сервиса необходим один бинарный файл.

## ⚡ Запуск сервиса

1. Скачать исполняемый
   файл [Linux](https://github.com/advancemg/vimb-loader/releases/download/untagged-1a1083e2fdbf15f8e9f9/vimb-loader-linux-amd64.zip)
   , [Windows](https://github.com/advancemg/vimb-loader/releases/download/untagged-1a1083e2fdbf15f8e9f9/vimb-loader-windows-amd64.zip)
   под Вашу ОС.
2. Распаковать скачанный zip файл, открыть терминал и выполниь команду:
   ``` bash 
   unzip Download/vimb-loader-linux-amd64.zip.
   ```
3. Дать права на выполнение файла, в терминале выполнить команду:
   ``` bash 
   sudo chmod 777 vimb-loader-linux-amd64
   ```
4. Запустить программу командой:
   ``` bash 
   ./vimb-loader-linux-amd64 -config
   ```
5. В терминале появится контекстное меню ``Edit config? (Y/n):``, при ответе ``Y`` можно будет отредактировать файл
   конфигурации ``config.json``.
6. Настройка конфигурации vimb-loader (при нажатии на Enter без ввода значения установится дефолтное значение):
    * Задайте переодичность скачивания бюджетов. ``Enter Budget cron(default 0 0/46 * * *):``
    * Укажите направление продаж. ``Enter Budget sellingDirection(default 23):``
    * Включить загрузку бюджетов по расписанию? ``Budget loading? (default false):``
    * Задайте переодичность скачивания сеток. ``Enter ProgramBreaks cron(default 0 0 0/8 * *):``
    * Укажите направление продаж. ``Enter ProgramBreaks sellingDirection(default 23):``
    * Включить загрузку сеток по расписанию? ``ProgramBreaks loading? (default false):``
    * Задайте переодичность скачивания сеток. ``Light Mode Enter ProgramBreaksLight cron(default 0/2 * * * *):``
    * Укажите направление продаж. ``Enter ProgramBreaksLight sellingDirection(default 23):``
    * Включить загрузку сеток. ``Light Mode по расписанию? ProgramBreaksLight loading? (default false):``
    * Задайте переодичность скачивания медиапланов. ``Enter Mediaplan cron(default 0 0/20 * * *):``
    * Укажите направление продаж. ``Enter Mediaplan sellingDirection(default 23):``
    * Включить загрузку медиапланов по расписанию? ``Mediaplan loading? (default false):``
    * Задайте переодичность скачивания списока роликов. ``Enter AdvMessages cron(default 0 0/2 * * *):``
    * Укажите направление продаж. ``Enter AdvMessages sellingDirection(default 23):``
    * Включить загрузку список роликов по расписанию? ``AdvMessages loading? (default false):``
    * Задайте переодичность скачивания спотов. ``Enter Spots cron(default 0 0/59 * * *):``
    * Укажите направление продаж. ``Enter Spots sellingDirection(default 23):``
    * Включить загрузку спотов по расписанию? ``Spots loading? (default false):``
    * Задайте переодичность скачивания удаленных за период
      спотов. ``Enter DeletedSpotInfo cron(default 0 0 0/12 * *):``
    * Укажите направление продаж. ``Enter DeletedSpotInfo sellingDirection(default 23):``
    * Включить загрузку удаленных за период спотах по расписанию? ``DeletedSpotInfo loading? (default false):``
    * Задайте переодичность скачивания каналов. ``Enter Channels cron(default 0 0/18 * * *):``
    * Укажите направление продаж. ``Enter Channels sellingDirection(default 23):``
    * Включить загрузку каналов по расписанию? ``Channels loading? (default false):``
    * Задайте переодичность скачивания списока заказчиков с
      рекламодеталями. ``Enter CustomersWithAdvertisers cron(
      default 0 0/16 * * *):``
    * Укажите направление продаж. ``Enter CustomersWithAdvertisers sellingDirection(default 23):``
    * Включить загрузку списока заказчиков с рекламодеталями по
      расписанию? ``CustomersWithAdvertisers loading? (default
      false):``

    * Задайте переодичность скачивания справочник рангов размещения. ``Enter Rank cron(default 0 0 0/23 * *):``
    * Включить загрузку справочника рангов размещения по расписанию? ``Rank loading? (default false):``
    * Введите адрес хоста RabbitMQ (если указать локальный хост запустится встроенный RabbitMQ). ``Enter amqp host(
      default localhost):``
    * Введите порт RabbitMQ. ``Enter amqp port(default 5555):``
    * Введите пользователя RabbitMQ. ``Enter amqp username(default guest):``
    * Введите пароль RabbitMQ. ``Enter amqp password(default guest):``
    * Введите пользователя MinIO. ``Enter S3 AccessKeyId(default minioadmin):``
    * Введите пароль MinIO. ``Enter S3 SecretAccessKey(default minioadmin):``
    * Введите адрес хоста MinIO (если указать локальный хост запустится встроенный MinIO). ``Enter S3 Endpoint(
      default 127.0.0.1:9999):``
    * Выберите режим соединения, true без шифрования. ``Enter S3 Debug(default true):``
    * Введите название бакета. ``Enter S3 Bucket(default storage):``
    * Укажите локальную директорию для хранения данных MinIO. ``Enter S3 LocalDir(default s3-buckets):``
    * Введите адрес ВИМБ сервиса. 435 боевой, 436
      тестовый. ``Enter url(default https://vimb-svc.vitpc.com:436/VIMBService.asmx):``
    * Введите сертификат выданный ВИМБ в формате base64. ``Enter cert:``
    * для преобразовая сертификата в base64 выволните команду:
    * "Certificate_file=путь к сертификату"
    ``` bash
     echo | base64 $Certificate_file | tr -d '\r\n'
    ```
    * Введите пароль выданный ВИМБ. ``Enter password:``
    * Введите клиента, используется при нейминге s3Key. ``Enter client(default test_client):``
    * Введите timeout на соединения с сервисом ВИМБ, обязателен суффикс времени s или
   m. ``Enter timeout(default 120s):``

7. Vimb-loader начнет работу согласно заданной конфигурации. По
   адрсу [localhost:8000/api/v1/docs/index.html](http://localhost:8000/api/v1/docs/index.html) можно ознакомиться с
   документацией API.

Периодичность указывается в формате cron 
```
* * * * * *
| | | | | |
| | | | | +--- Годы       (диапазон: 1900-3000)
| | | | +----- Дни недели (диапазон: 1-7)
| | | +------- Месяцы     (диапазон: 1-12)
| | +--------- Дни месяца (диапазон: 1-31)
| +----------- Часы       (диапазон: 0-23)
+------------- Минуты     (диапазон: 0-59)
```