package mst

import (
    "io"
    "os"
    "log"
    "os/user"
    "net/http"
    "encoding/json"
    "github.com/rakyll/statik/fs"
    "github.com/ileler/hybrid-launcher"
    // Import ./static/*
    _ "github.com/ileler/miio-serial-tools/statik"
)

var logger *log.Logger
var configFile string
var config struct {
    Port int
}

func Main() {

    http.HandleFunc("/exit", func (w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        launcher.Exit()
    })

    statikFS, err := fs.New()
    if err != nil {
        panic(err)
    }
    http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
    //http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

    myself, error := user.Current()
    if error != nil {
        panic(error)
    }
    homedir := myself.HomeDir + "/.mst/"
    if err := os.MkdirAll(homedir, 0775); err != nil {
        panic(err)
    }

    lf, err := os.OpenFile(homedir + "log.txt", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
    if err != nil {
        panic(err)
    }
    //defer lf.Close()

    logger = log.New(io.MultiWriter(lf, os.Stdout), "", log.LstdFlags)
    //logger.SetOutput(os.Stdout)

    api(homedir)

    configFile, err := os.Open(homedir + ".cfg")
    defer configFile.Close()
    port := 0
    if err == nil {
        jsonParser := json.NewDecoder(configFile)
        if err = jsonParser.Decode(&config); err == nil {
            logger.Println("config: ", config)
            port = config.Port
        } else {
            logger.Println("decode: ", err)
        }
    } else {
        logger.Println("open: ", err)
    }

    //c := &launcher.Config{ Pid: homedir + ".pid", Port: port, HandleRoot: false }
    c := launcher.DefaultConfig()
    c.Pid = homedir + ".pid"
    c.Port = port
    c.HandleRoot = false
    launcher.StartWithConfig(c)
    //launcher.Start()

}
