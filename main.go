
    package main

    import (
    	"flag"
    	"fmt"
    	"github.com/confluentinc/confluent-kafka-go/kafka"
    	"strconv"
    	"sync"
    	"time"
    )

    func main() {
    	//file, err := os.ReadFile("file.png")
    	//if err != nil {
    	//	fmt.Println(err)
    	//	return
    	//
    	//}

    	var serverName = ""
    	flag.StringVar(&serverName, "server-name", "default", "help diagnosing witch application produce messages to kafka")

    	var bootstrapServers = ""
    	flag.StringVar(&bootstrapServers, "bootstrap-servers", "", "bootstrap servers of kafka")

    	var topic = ""
    	flag.StringVar(&topic, "topic", "topic", "kafka topic that you want to publish to")

    	var batchCount = 0
    	flag.IntVar(&batchCount, "batch-count", 0, "number of message batches")

    	var batchSize = 0
    	flag.IntVar(&batchSize, "batch-size", 0, "count of message that you want to concurrently publish to kafka")

    	flag.Parse()

    	fmt.Println( "serverName:", serverName)
    	fmt.Println( "bootstrapServers:", bootstrapServers)
    	fmt.Println( "topic:", topic)
    	fmt.Println( "batchCount:", batchCount)
    	fmt.Println( "batchSize:", batchSize)
    	fmt.Println()

    	fmt.Println( "create kafka producer")
    	newProducer, err := kafka.NewProducer(&kafka.ConfigMap{
    		"bootstrap.servers": bootstrapServers,
    	})
    	if err != nil {
    		fmt.Println(err)
    		return
    	}

    	fmt.Println( "start producing message in kafka")
    	start := time.Now()
    	for i := 0; i < batchCount; i++ {
    		//s := time.Now()
    		produce(newProducer,nil, serverName, i, batchSize, topic)
    		//fmt.Println( serverName, " batch#"+strconv.Itoa(i)+" published to kafka", "in", time.Since(s))
    	}
    	fmt.Println()
    	fmt.Println( serverName, "publish", batchCount*batchSize, "messages in ", time.Since(start))
    	fmt.Println( serverName, "done!")
    }

    func produce(producer *kafka.Producer, file []byte, serverName string, batchNumber int, count int, topic string) {
    	var wg sync.WaitGroup
    	wg.Add(count)

    	for i := 0; i < count; i++ {

    		go func(i int) {
    			defer wg.Done()
    			eventChan := make(chan kafka.Event)
    			err := producer.Produce(&kafka.Message{
    				TopicPartition: kafka.TopicPartition{
    					Topic: &topic,
    					//Partition: int32(i % partitionCount),
    					Partition: kafka.PartitionAny,
    				},
    				Value:         []byte(serverName + "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq message#" + strconv.Itoa(i) + " in batch#" + strconv.Itoa(batchNumber)),
    				//Value:         file,
    				Key:           nil,
    				Timestamp:     time.Time{},
    				TimestampType: 0,
    				Opaque:        nil,
    				Headers:       nil,
    			}, eventChan)
    			<-eventChan
    			if err != nil {
    				fmt.Println(serverName, "publish error", err)
    			}
    		}(i)
    	}
    	wg.Wait()
    }
