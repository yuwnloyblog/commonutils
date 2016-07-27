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
	if self.Props == nil {
		return "",&MissingError{"Please load properties file before use."}
	}
	if val,ok := self.Props[key];ok{
		return val,nil
	}
	return "",&MissingError{"Missing key,"+key}
}

func (self *Properties)GetStringWithDefault(key,defaultVal string)string{
	val,err := self.GetString(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (self *Properties)GetInt(key string)(int,error){
	if self.Props == nil {
		return 0,&MissingError{"Please load properties file before use."}
	}
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
func (self *Properties)GetIntWithDefault(key string, defaultVal int)int{
	val,err := self.GetInt(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (self *Properties)GetBool(key string)(bool,error){
	if self.Props == nil {
		return false,&MissingError{"Please load properties file before use."}
	}
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
func (self *Properties)GetBoolWithDefault(key string, defaultVal bool)bool{
	val,err := self.GetBool(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (self *Properties)GetFloat(key string)(float64,error){
	if self.Props == nil {
		return float64(0),&MissingError{"Please load properties file before use."}
	}
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

func (self *Properties)GetFloatWithDefault(key string, defaultVal float64)float64{
	val,err := self.GetFloat(key)
	if err != nil {
		return defaultVal
	}
	return val
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