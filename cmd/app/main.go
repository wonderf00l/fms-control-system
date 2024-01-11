package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/wonderf00l/fms-control-system/internal/app"
)

// Connect to the broker and publish a message periodically

const (
	TOPIC         = "topic1"
	QOS           = 1
	SERVERADDRESS = "tls://mqtt.cloud.yandex.net:8883"
	DELAY         = time.Second
	CLIENTID      = "mqtt_publisher"

	WRITETOLOG = true // If true then published messages will be written to the console
)

func NewTLSConfig() *tls.Config {
	rootCAs := x509.NewCertPool()
	clientCAs := x509.NewCertPool()

	brokerCert, err := os.ReadFile("../../certs/rootCA.crt")
	if err != nil {
		log.Fatalln(err)
	}
	clientCert, err := os.ReadFile("../../certs/self-signed.crt")
	if err != nil {
		log.Fatalln(err)
	}

	rootCAs.AppendCertsFromPEM(brokerCert)
	clientCAs.AppendCertsFromPEM(clientCert)

	cert, err := tls.LoadX509KeyPair("../../certs/client.crt", "../../certs/client.key")
	if err != nil {
		log.Fatalln(err)
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Fatalln(err)
	}

	return &tls.Config{
		RootCAs:      rootCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
		Certificates: []tls.Certificate{cert},
	}

}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
}

// func main() {
// 	// Enable logging by uncommenting the below
// 	// mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
// 	// mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
// 	// mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
// 	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
// 	opts := mqtt.NewClientOptions()
// 	opts.SetTLSConfig(NewTLSConfig())
// 	opts.SetCleanSession(false)
// 	opts.AddBroker(SERVERADDRESS)
// 	opts.SetClientID(CLIENTID)
// 	opts.SetUsername("a44p9bq63uah16qnocbi")
// 	opts.SetPassword("iot-course-worker22036")
// 	cwd, _ := os.Getwd()
// 	fmt.Println(cwd)
// 	opts.SetOrderMatters(false)           // Allow out of order messages (use this option unless in order delivery is essential)
// 	opts.ConnectTimeout = 5 * time.Second // Minimal delays on connect
// 	opts.WriteTimeout = 5 * time.Second   // Minimal delays on writes
// 	opts.KeepAlive = 50                   // Keepalive every 10 seconds so we quickly detect network outages
// 	opts.PingTimeout = 5 * time.Second    // local broker so response should be quick

// 	// Automate connection management (will keep trying to connect and will reconnect if network drops)
// 	opts.ConnectRetry = true
// 	opts.AutoReconnect = true

// 	// Log events
// 	opts.OnConnectionLost = func(cl mqtt.Client, err error) {
// 		fmt.Println("connection lost")
// 	}
// 	opts.OnConnect = func(c mqtt.Client) {
// 		if token := c.Subscribe(TOPIC, QOS, onMessageReceived); token.Wait() && token.Error() != nil {
// 			panic(token.Error())
// 		}
// 	}
// 	opts.OnReconnecting = func(mqtt.Client, *mqtt.ClientOptions) {
// 		fmt.Println("attempting to reconnect")
// 	}
// 	opts.SetDefaultPublishHandler()
// 	//
// 	// Connect to the broker
// 	//
// 	client := mqtt.NewClient(opts)

// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		panic(token.Error())
// 	}
// 	defer client.Disconnect(500)

// 	time.Sleep(5 * time.Second)
// 	client.Publish(TOPIC, QOS, false, "test")
// 	fmt.Println("Connection is up")

// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

// 	<-sig
// }

func main() {
	serviceLogger, cfgFiles, err := initPrerequisites()
	if err != nil {
		log.Fatal(err)
	}
	defer serviceLogger.Sync()

	if err := app.Run(serviceLogger, cfgFiles); err != nil {
		serviceLogger.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	serviceLogger.Infoln("Shutting down gracefully")
}
