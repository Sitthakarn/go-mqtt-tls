package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"

    MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewTLSConfig() *tls.Config {
    // CA証明書を設定
    certpool := x509.NewCertPool()
    pemCerts, err := ioutil.ReadFile("xxxxxxxxx.pem")
    if err == nil {
        certpool.AppendCertsFromPEM(pemCerts)
    }

    // クライアント証明書とキーペアを設定
    cert, err := tls.LoadX509KeyPair("xxxxxxxxx-certificate.pem.crt", "xxxxxxxxx-private.pem.key")
    if err != nil {
        panic(err)
    }

    cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
    if err != nil {
        panic(err)
    }

    // config設定
    return &tls.Config{

        RootCAs: certpool,
        ClientAuth: tls.NoClientCert,
        ClientCAs: nil,
        InsecureSkipVerify: true,
        Certificates: []tls.Certificate{cert},
    }
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
 	   fmt.Printf("TOPIC: %s\n", msg.Topic())
    fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
    tlsconfig := NewTLSConfig()

    opts := MQTT.NewClientOptions()
    opts.AddBroker("ssl://xxxxxxxxxxxxxxxxxxx.amazonaws.com:8883")
    opts.SetClientID("ssl-sample").SetTLSConfig(tlsconfig)
    opts.SetDefaultPublishHandler(f)

    // 接続をする
    c := MQTT.NewClient(opts)
    if token := c.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    fmt.Println("AWS IoT Connect Success")

    if token := c.Publish("$aws/things/ESP32/shadow/update", 0, false,
        `{"state":{"reported":{"welcome":"I am gopher!!!"}}}`); token.Wait() && token.Error() != nil {
                panic(token.Error())
    }
    fmt.Println("Message Publish Success")

    // 切断
    c.Disconnect(250)
}
