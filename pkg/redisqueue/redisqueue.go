

func ( p *Producer ) OpenProducer() { 
	p, err := redisqueue.NewProducerWithOptions(&redisqueue.ProducerOptions{
		StreamMaxLength:      10000,
		ApproximateMaxLength: true,
		RedisOptions.Network: localhost:6379,
	})
	if err != nil {
		panic(err)
	}
}

func ( p *Producer ) ProduceMessage( queueName string) { 
		err := p.Enqueue(&redisqueue.Message{
			Stream: queueName,
			Values: map[string]interface{}{
				"index": i,
			},
		})
		if err != nil {
			panic(err)
		}
}

func (c *Consumer ) OpenConsumer( queueName string ) {

	c, err := redisqueue.NewConsumerWithOptions(&redisqueue.ConsumerOptions{
		Name: t.Format("20060102150405"),
		VisibilityTimeout: 60 * time.Second,
		BlockingTimeout:   5 * time.Second,
		ReclaimInterval:   1 * time.Second,
		BufferSize:        100,
		Concurrency:       10,
		RedisOptions.Network: localhost:6379,		
	})
	if err != nil {
		panic(err)
	}
	
	c.Register(queueName, process)
}


func (c *Consumer ) ConsumerEnqueue( queueName ) {

		err := p.Enqueue(&redisqueue.Message{
			Stream: queueName,
			Values: map[string]interface{}{
				"index": i,
			},
		})
		if err != nil {
			panic(err)
		}
}

