package main
import (
    //"encoding/json"
    //"flag"

    "fmt"
    "github.com/ChimeraCoder/anaconda"
    "time"
    "github.com/claudiu/gocron"
    _ "github.com/ziutek/mymysql/godrv"
    "database/sql"
    "encoding/json"
    "os"
    "bytes"
    "net/http"
    "io/ioutil"
    "strings"
)

var con *sql.DB
func check(e error) {
    if e != nil {
        panic(e)
    }
}

var database string
var user string
var password string

type Configuration struct {
    DB    string
    US   string
    PASS string
}

type Usuario struct{

  idUsuario int
  idRedSocial string
  nombreUsuarioRedSocial string
  nombreU string
  apellU string
  key string
}

func conexionOn(){
    conn,err:= sql.Open("mymysql",database+"/"+user+"/"+password)
    if err!=nil{
      fmt.Println("No Conectado a BD")
    }
    //asigno la var local de conexion a la global
    con=conn
    return
}

func main(){


  fmt.Println("**********************************************************")
  fmt.Println("                 Inicio Motor TWITTER                         ")
  fmt.Println("**********************************************************")

  setCredenciales()

  s := gocron.NewScheduler()
  s.Every(300).Seconds().Do(OrquestadorFlujo)
  //s.Every(100).Seconds().Do(OrquestadorFlujo)
  sc := s.Start() // keep the channel
  //go test(s, sc)  // wait
  <-sc            // it will happens if the channel is closed*/

}

func test(s *gocron.Scheduler, sc chan bool) {
    time.Sleep(8 * time.Second)
    s.Clear()
    fmt.Println("All task removed")
    close(sc) // close the channel
}

func setCredenciales(){

  //Vamos a llamar un archivo de configuraciòn para obtener las credenciales de la bd
  file, _ := os.Open("conf.json")
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)
  if err != nil {
  fmt.Println("error:", err)
  }
  //fmt.Println(configuration) // output: [UserA, UserB]
  user=configuration.US
  password=configuration.PASS
  database=configuration.DB

}

func getDatosUsuarios() []Usuario{

    row,err:=con.Query("select r.id_usuario, r.id_usuario_red,r.n_usuario_red,u.keyMovil,u.nom1,u.apell1 from rs_usuario r, usuario u where n_usuario_red<>'sin_info' and id_red_social=2 and r.id_usuario = u.id_usuario")
    //row := con.Query("select id_usuario, id_usuario_red,n_usuario_red from rs_usuario where n_usuario_red<>'sin_info' and id_red_social=1 and id_usuario in( select id_usuario from usuario)")
    if err != nil {
      fmt.Println("Error de Conexion")
    }

    var Usuarios []Usuario
    var idUsuario int
    var idRedSocial string
    var nombreUsuarioRedSocial string
    var key string
    var nombreU string
    var apellU string

    for row.Next(){
      row.Scan(&idUsuario,&idRedSocial,&nombreUsuarioRedSocial,&key,&nombreU,&apellU)
      //fmt.Println(nombreUsuarioRedSocial)
      Usuarios= append(Usuarios,Usuario{idUsuario,idRedSocial,nombreUsuarioRedSocial,key,nombreU,apellU} )
    }

    return Usuarios
}

func aString(s []string) string{
  var buffer bytes.Buffer

    for i := 0; i < len(s); i++ {
        buffer.WriteString(s[i])
    }

    //fmt.Println("esta es la query:"+buffer.String())
    return buffer.String()
}

func rutinaGo(chanel []Usuario ){
  //var sqlstring string
  for cont:=0 ;  cont< len(chanel); cont++{
        //fmt.Println(chanel[cont].idUsuario)
        ListaPalabras:=getPalabrasMalasPorUsuario(chanel[cont].idUsuario)
        construyeInserts(ListaPalabras,chanel[cont].nombreUsuarioRedSocial,chanel[cont].idUsuario,chanel[cont].key,chanel[cont].nombreU,chanel[cont].apellU)

    }
}

func push( keymovil string){
  searchURL := "http://192.168.30.206:7010/fimi_v0/webapi/u/SPush"
  tosearch:=searchURL+";id="+ keymovil+";cod=1;contenido=desde twitter"
  response,err := http.Get(tosearch)
  if err != nil {
      fmt.Printf("%s", err)
      os.Exit(1)
  } else {
      defer response.Body.Close()
      contents, err := ioutil.ReadAll(response.Body)
      if err != nil {
          fmt.Printf("%s", err)
          os.Exit(1)
      }
      fmt.Printf("%s\n", string(contents))
  }
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func construyeInserts(palabras []string, nom string,idUsuario int, keymovil string,nombreU string,apellU string){
  var i int
   t := time.Now()
  fmt.Println("**********************************************************")
  fmt.Println("                 LLAMANDO TWITTER                         ")
  fmt.Println("**********************************************************")
  fmt.Println(nom)
  fmt.Println(palabras)
  fmt.Println("Fecha Actual de servicio:",t)
  fmt.Println("**********************************************************")

  i= i + rutina1(palabras,nom,idUsuario)
  i= i + rutina2(palabras,idUsuario,nombreU,apellU)
  i= i + rutina3(palabras,idUsuario,nom)

  if i > 0{
    push(keymovil)
  }
  fmt.Println("**********************************************************")
  fmt.Println("                 Fin Proceso Twiiter                        ")
  fmt.Println("**********************************************************")
}

func rutina3(palabras []string,idUsuario int,nom string)int{
  fmt.Println("Rutina 3!")
  var i int
  var cant int
  for conta:=0; conta < len(palabras); conta++{
    cant=1
    //fmt.Println(palabras[conta])
    //trabajo api twitter
    anaconda.SetConsumerKey("bCsERW0NYC129jSwX5kaNbXRD")
    anaconda.SetConsumerSecret("NPVbMA5pWhAV7By7ORML6kTuhUv1ls07wYZgC8OLVR4kMyo3Dk")
    api := anaconda.NewTwitterApi("56242109-6PGhfzlLkMc1TZTFKCAuVfdMBLTY7dlso7EH1IpYk", "TN1r6fSTMGf7112EyPBeQeAyUETS7Xg7qxglGGNorCCLN")
    busqueda:=nom
    twitterStream, _:= api.GetSearch(busqueda, nil)
    /*if (idUsuario==13){
      fmt.Println(twitterStream.Statuses)
    }*/
    var texto string

    for _ , tweet := range twitterStream.Statuses {
        fmt.Println("-----------")
        fmt.Println(cant)
         fmt.Println(tweet.Text)
         texto=strings.Replace(tweet.Text,"'"," ",-1)
         //fmt.Println(texto)
         res,err:=con.Exec("INSERT into historial_usuario (id_usuario_red,id_red_social,comentario,tipo_comentario,fecha,id_onombre_quien_comenta,is_falso_positivo) values ("+fmt.Sprintf("%d",idUsuario)+",2,'"+texto+"','-',sysdate(),'sin_info',1)")
         if err != nil {
             fmt.Println("")
             fmt.Println("insert ERROR: Twit ya existe")
             fmt.Println("")
         }else{
           id,err:= res.LastInsertId()
           if err != nil {
                fmt.Println("")
                fmt.Println("insert ERROR: Twit ya existe")
                fmt.Println("")
           }
           if(id==0){
             i=1
             fmt.Println("")
             fmt.Println("insert OK")
             fmt.Println("")
           }
         }
      cant++
    }

  }
  fmt.Println("FIN Rutina 3!")
  return i
}

func rutina2(palabras []string,idUsuario int,nombreU string,apellU string) int{
  fmt.Println("Rutina 2!")
  var i int
  var cant int
  for conta:=0; conta < len(palabras); conta++{
    cant=1
    //fmt.Println(palabras[conta])
    //trabajo api twitter
    anaconda.SetConsumerKey("bCsERW0NYC129jSwX5kaNbXRD")
    anaconda.SetConsumerSecret("NPVbMA5pWhAV7By7ORML6kTuhUv1ls07wYZgC8OLVR4kMyo3Dk")
    api := anaconda.NewTwitterApi("56242109-6PGhfzlLkMc1TZTFKCAuVfdMBLTY7dlso7EH1IpYk", "TN1r6fSTMGf7112EyPBeQeAyUETS7Xg7qxglGGNorCCLN")
    busqueda:=nombreU+" "+apellU+" "+palabras[conta]
    twitterStream, _:= api.GetSearch(busqueda, nil)
    /*if (idUsuario==13){
      fmt.Println(twitterStream.Statuses)
    }*/
    var texto string

    for _ , tweet := range twitterStream.Statuses {
        fmt.Println("-----------")
        fmt.Println(cant)
         fmt.Println(tweet.Text)
         texto=strings.Replace(tweet.Text,"'"," ",-1)
         //fmt.Println(texto)
         res,err:=con.Exec("INSERT into historial_usuario (id_usuario_red,id_red_social,comentario,tipo_comentario,fecha,id_onombre_quien_comenta,is_falso_positivo) values ("+fmt.Sprintf("%d",idUsuario)+",2,'"+texto+"','negativo',sysdate(),'sin_info',1)")
         if err != nil {
             fmt.Println("")
             fmt.Println("insert ERROR: Twit ya existe")
             fmt.Println("")
         }else{
           id,err:= res.LastInsertId()
           if err != nil {
                fmt.Println("")
                fmt.Println("insert ERROR: Twit ya existe")
                fmt.Println("")
           }
           if(id==0){
             i=1
             fmt.Println("")
             fmt.Println("insert OK")
             fmt.Println("")
           }
         }
      cant++
    }

  }
  fmt.Println("FIN Rutina 2!")
  return i
}
func rutina1(palabras []string, nom string,idUsuario int) int {
    fmt.Println("Rutina 1!")
    var i int
    var cant int
    for conta:=0; conta < len(palabras); conta++{
      cant=1
      //fmt.Println(palabras[conta])
      //trabajo api twitter
      anaconda.SetConsumerKey("bCsERW0NYC129jSwX5kaNbXRD")
      anaconda.SetConsumerSecret("NPVbMA5pWhAV7By7ORML6kTuhUv1ls07wYZgC8OLVR4kMyo3Dk")
      api := anaconda.NewTwitterApi("56242109-6PGhfzlLkMc1TZTFKCAuVfdMBLTY7dlso7EH1IpYk", "TN1r6fSTMGf7112EyPBeQeAyUETS7Xg7qxglGGNorCCLN")
      busqueda:=nom+" "+palabras[conta]
      twitterStream, _:= api.GetSearch(busqueda, nil)
      /*if (idUsuario==13){
        fmt.Println(twitterStream.Statuses)
      }*/
      var texto string

      for _ , tweet := range twitterStream.Statuses {
          fmt.Println("-----------")
          fmt.Println(cant)
           fmt.Println(tweet.Text)
           texto=strings.Replace(tweet.Text,"'"," ",-1)
           //fmt.Println(texto)
           res,err:=con.Exec("INSERT into historial_usuario (id_usuario_red,id_red_social,comentario,tipo_comentario,fecha,id_onombre_quien_comenta,is_falso_positivo) values ("+fmt.Sprintf("%d",idUsuario)+",2,'"+texto+"','negativo',sysdate(),'sin_info',1)")
           if err != nil {
               fmt.Println("")
               fmt.Println("insert ERROR: Twit ya existe")
               fmt.Println("")
           }else{
             id,err:= res.LastInsertId()
             if err != nil {
                  fmt.Println("")
                  fmt.Println("insert ERROR: Twit ya existe")
                  fmt.Println("")
             }
             if(id==0){
               i=1
               fmt.Println("")
               fmt.Println("insert OK")
               fmt.Println("")
             }
           }
        cant++
      }

    }
    fmt.Println("FIN Rutina 1!")
    return i
}

func OrquestadorFlujo(){
  conexionOn()
  var grupo1 []Usuario
  var grupo2 []Usuario
  var grupo3 []Usuario
  var e int
  ListaUsuarios:=getDatosUsuarios()
  fmt.Println(ListaUsuarios)
  e=0
  for cont:=0 ;  cont< len(ListaUsuarios); cont++{
    if e==0{
      grupo1=append(grupo1,Usuario{ListaUsuarios[cont].idUsuario,ListaUsuarios[cont].idRedSocial,ListaUsuarios[cont].nombreUsuarioRedSocial,ListaUsuarios[cont].key,ListaUsuarios[cont].nombreU,ListaUsuarios[cont].apellU} )
    }else if e==1{
            grupo2=append(grupo2,Usuario{ListaUsuarios[cont].idUsuario,ListaUsuarios[cont].idRedSocial,ListaUsuarios[cont].nombreUsuarioRedSocial,ListaUsuarios[cont].key,ListaUsuarios[cont].nombreU,ListaUsuarios[cont].apellU} )
          }else {
              grupo3=append(grupo3,Usuario{ListaUsuarios[cont].idUsuario,ListaUsuarios[cont].idRedSocial,ListaUsuarios[cont].nombreUsuarioRedSocial,ListaUsuarios[cont].key,ListaUsuarios[cont].nombreU,ListaUsuarios[cont].apellU} )
              e=-1
          }
    e++
  }

  go rutinaGo(grupo1)
  go rutinaGo(grupo2)
  go rutinaGo(grupo3)
}

func getPalabrasMalasPorUsuario(idUsuario int) []string{

  sQuery:="select palabra from palabras_usuario where tipo_palabra=1 and id_usuario="+fmt.Sprintf("%d",idUsuario)
  row,err:=con.Query(sQuery)
  if err!=nil{
    fmt.Println("Error de Conexion")
   }
   var palabra string
   var ListaPalabras []string

   for row.Next(){
     row.Scan(&palabra)
     ListaPalabras=append(ListaPalabras,palabra)
   }

   sQuery2:="select palabra from palabras_fimi where tipo_palabra=1"
   row2,err:=con.Query(sQuery2)
   if err!=nil{
     fmt.Println("Error de Conexion")
   }
   for row2.Next(){
     row2.Scan(&palabra)
     ListaPalabras=append(ListaPalabras,palabra)
   }



   return ListaPalabras
}
