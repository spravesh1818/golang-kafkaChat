# CLI chatroom Application Using Kafka As Message Queue

This app demonstrates how we can use kafka as backend instead of popular websockets. It leverages kafka to produce and consume messages and show them on the cli.Here are the steps to install and use the application.

### Pre-Installation

* First of all you will need a working copy of kafka locally. For that you will need to use the docker-compose.yml file.

  ```
  docker-compose up -d
  ```

  This will pull the kafka image and run the containers.Kafka will be deployed at localhost:9092 while zookeeper runs on localhost:2181
* Next you will need to create the topic.For this I have included the create_topic.sh file.We will need to make it executable and run it through the command line.

  ```
  chmod +x create_topic.sh
  ```

  Once done making the file executable just use the following command to create topic.

  ```
  ./create_topic <topic_name>
  ```

  The <topic_name> here is very important as this topic name will be used in our env file as the base topic where we produce and consume messages
* Once this is done we will need to add values to our env file.Copy the .env.sample file and then update the topic_name created in the above steps.

  ```
  mv .env.example .env
  ```

  *Note: If topic is not present we cannot continue as kafka primarily uses topic to publish or subscribe to.*

### Installation

1. First clone the repository.

   ```
   git clone git@github.com:spravesh1818/golang-kafkaChat.git
   ```
2. Install the dependencies.

   ```
   go install
   ```
3. Run one instance of server.I have included a Makefile so that it is easy to run.

   ```
   make server-run
   ```
4. Run client.Client is where you can run send message and see messages.Since it is a cli application still it is not that user-friendly.But remember this is just an application to show the capability of kafka as a message queue.

   ```
   make client-run
   ```

   *Note:You can run multiple instances of client and interact as two users to chat with each other.*

### Technologies and Frameworks Used

* **Language:** Golang
* **Frameworks:** Gin
* **Message Broker:** Apache Kafka
* **Shell Scripting :** Bash
