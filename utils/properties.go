package utils
import(
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
	"strings"
	"strconv"
)

type MissingError struct{
	What string
}
func (self *MissingError) Error() string {
	return fmt.Sprintf("at %v, %s",
		time.Now(), self.What)
}

type Properties struct{
	Props map[string]string
}

func (self *Properties)GetString(key string)(string,error){
	if val,ok := self.Props[key];ok{
		return val,nil
	}
	return "",&MissingError{"Missing key,"+key}
}

func (self *Properties)GetInt(key string)(int,error){
	if val,ok := self.Props[key];ok{
		intVal,err:=strconv.Atoi(val)
		if err == nil {
			return intVal,nil
		}else{
			return 0,err
		}
	}
	return 0,&MissingError{"Missing key,"+key}
}

func (self *Properties)GetBool(key string)(bool,error){
	if val,ok := self.Props[key];ok{
		bVal,err := strconv.ParseBool(val)
		if err == nil{
			return bVal,nil
		}else{
			return false,err
		}
	}
	return false,&MissingError{"Missing key,"+key}
}

func (self *Properties)GetFloat(key string)(float64,error){
	if val,ok := self.Props[key];ok{
		fVal,err := strconv.ParseFloat(val,64)
		if err == nil{
			return fVal,nil
		}else{
			return float64(0),err
		}
	}
	return float64(0),&MissingError{"Missing key,"+key}
}

/**
 * read the properties file
 */
func (self *Properties)LoadPropertyFile(filepath string)(error){
	props := make(map[string]string)
	file,err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		bytes,_,err := reader.ReadLine()
		if err != nil {
			if err != io.EOF{
				return err
			}else{
				break
			}
		}
		line := string(bytes)
		isContainedSplit := strings.Contains(line,"=")
		if isContainedSplit {
			arr := strings.Split(line,"=")
			if len(arr) >= 2 {
				props[arr[0]] = arr[1]
			}
		}
	}
	self.Props = props
	return nil
}