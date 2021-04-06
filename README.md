# Messaging Systems Comparison
This repo is prepared to test the performance of the following messaging systems:

- Kafka
- Nats
- RabbitMQ
- Redis

For each messaging system this repo measures the following properties:

- Total time of sending
- Total time of receiving (starting from first receiving)
- Total time from first sending and last receiving
- Latency: time spent to publish message
- Throughput: number of published messages by second
- Transmission time: time spent in the messaging system (from sending to its receiving)
- Bandwidth: number of messages that are received by second

## PREQUISITES

- Golang >1.16
- Docker
- Docker compose

## FOLDERS
```
.
├── benchmark                          Benchmark go struct
├── benchmarks                         Legacy benchmark codes
├── cmd/messagging-system-benchmark    Main command code
├── delivery                           Deployment folder
├── faker                              Faker command code
├── json                               Mockups folder
├── logs                               Logs folder (Unused)
├── measurement                        Measurement struct code
├── model                              Message model
├── services                           Messaging services code
├── docker-compose.yaml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```
### faker
With faker, the user can generate json files with fake messages. Usage:
```sh
// In faker folder (cd faker)
$ go run faker.go
```
### json
JSON files generated with faker are located in this folder. Default files and sizes are listed below:
```shell
// In json folder (cd json)
$ ls -lah
total 113M
drwxrwxr-x  2 ijtr ijtr 4,0K mar 30 14:21 .
drwxrwxr-x 15 ijtr ijtr 4,0K abr  6 14:08 ..
-rw-rw-r--  1 ijtr ijtr 102M mar 30 14:22 100k.json
-rw-rw-r--  1 ijtr ijtr  10K mar 30 14:22 10.json
-rw-rw-r--  1 ijtr ijtr  11M mar 30 14:22 10k.json
-rw-rw-r--  1 ijtr ijtr 1,2K mar 30 14:22 1.json
-rw-rw-r--  1 ijtr ijtr 1,1M mar 30 14:22 1k.json
-rw-rw-r--  1 ijtr ijtr    0 mar 15 16:47 .gitkeep
```

### cmd/messagging-system-benchmark
This folder contains main file of the command to launch the benchmarks. To run the benchmark run the following command:
```sh
// In cmd/messagging-system-benchmark folder (cd cmd/messagging-system-benchmark)
$ go run main.go
```

### benchmark
This folder contains the struct for benchmark. The Benchmark struct represents a run with a file, a size of data and a messaging service. When a benchmark is run, it launch firstly the receiver and once it is ready, it launch the sender.

### model
This folder contains the struct for message. The Message struct has faker tags to generate meaningful dummy content like IP address etc.

### measurement
This folder contains the struct of the measurements of each benchmark. The Measurements struct consist of a collection of time measurements for each experiment including:
- total time
- total time between first sending and last receiving
- stats for each individual sending/receiving times

### services
This folder contains the main interface of a messaging service and all the implementations for each of the messaging systems. MessagingService interface consist of the use case of a messaging system with publish and subscription functions.

### benchmarks
These folders have same structure. They all have a sender and a receiver. The test steps are below. The order is important!

1. Start the relevant server according to the test (Kafka Streams, RabbitMQ Redis, Nats ~~or Nats Streaming Server~~).
2. Go into receiver folder of the test and run receiver with `go run receiver.go`
3. Go into sender folder of the test and run sender with `go run sender.go`

## RESULTS

All clients have been executed in a:

- Dell Latitude 7390 with OS Xubuntu 18.04, CPU Intel Core i5-8350U x4 and 16GiB DDR4 RAM Memory;
- using Docker v19.03.6 and Docker-Compose v1.25.4 for environment deployment; and
- Go 1.16.2 for clients execution.

The code bellow shows the output of the benchmark:

```
2021/04/06 14:13:53 ****Running Kafka benchmark****
2021/04/06 14:13:53 listening
2021/04/06 14:13:58 retreive last
2021/04/06 14:13:59 Processor shutdown cleanly
2021/04/06 14:13:59 Sender took: 4.3316 s
2021/04/06 14:13:59 Messages sent: 100000
2021/04/06 14:13:59 Latency: mean 1.4803e-05 s, variance 6.5732e-08 s
2021/04/06 14:13:59 Throughput: 67554 msg/s
2021/04/06 14:13:59 -------------------
2021/04/06 14:13:59 Receiver took: 2.2788 s
2021/04/06 14:13:59 Time among first send and last receive: 5.1629 s
2021/04/06 14:13:59 Messages received: 100000
2021/04/06 14:13:59 Transmission time: mean 1.888 s, variance 0.32968 s
2021/04/06 14:13:59 Bandwidth: 0.52966 msg/s
2021/04/06 14:13:59 ****************************
2021/04/06 14:13:59 ****Running Nats benchmark****
2021/04/06 14:13:59 listening
2021/04/06 14:14:02 Sender took: 3.6729 s
2021/04/06 14:14:02 Messages sent: 100000
2021/04/06 14:14:02 Latency: mean 1.1969e-05 s, variance 9.1929e-10 s
2021/04/06 14:14:02 Throughput: 83548 msg/s
2021/04/06 14:14:02 -------------------
2021/04/06 14:14:02 Receiver took: 3.6719 s
2021/04/06 14:14:02 Time among first send and last receive: 3.673 s
2021/04/06 14:14:02 Messages received: 100000
2021/04/06 14:14:02 Transmission time: mean 0.0003813 s, variance 3.396e-07 s
2021/04/06 14:14:02 Bandwidth: 2622.6 msg/s
2021/04/06 14:14:02 ****************************
2021/04/06 14:14:02 ****Running RabbitMq benchmark****
2021/04/06 14:14:02 listening
2021/04/06 14:14:09 Sender took: 6.7821 s
2021/04/06 14:14:09 Messages sent: 100000
2021/04/06 14:14:09 Latency: mean 4.0363e-05 s, variance 9.1879e-09 s
2021/04/06 14:14:09 Throughput: 24775 msg/s
2021/04/06 14:14:09 -------------------
2021/04/06 14:14:09 Receiver took: 6.7831 s
2021/04/06 14:14:09 Time among first send and last receive: 6.784 s
2021/04/06 14:14:09 Messages received: 100000
2021/04/06 14:14:09 Transmission time: mean 0.0026915 s, variance 1.0581e-05 s
2021/04/06 14:14:09 Bandwidth: 371.54 msg/s
2021/04/06 14:14:09 ****************************
2021/04/06 14:14:09 ****Running Redis benchmark****
2021/04/06 14:14:09 listening
2021/04/06 14:14:22 Sender took: 12.899 s
2021/04/06 14:14:22 Messages sent: 100000
2021/04/06 14:14:22 Latency: mean 0.00010045 s, variance 9.7132e-09 s
2021/04/06 14:14:22 Throughput: 9955.6 msg/s
2021/04/06 14:14:22 -------------------
2021/04/06 14:14:22 Receiver took: 12.899 s
2021/04/06 14:14:22 Time among first send and last receive: 12.899 s
2021/04/06 14:14:22 Messages received: 100000
2021/04/06 14:14:22 Transmission time: mean 0.00020111 s, variance 4.5721e-08 s
2021/04/06 14:14:22 Bandwidth: 4972.5 msg/s
2021/04/06 14:14:22 ****************************
```

