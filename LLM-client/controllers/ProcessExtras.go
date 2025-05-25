package controllers

import (
	"encoding/base64"
	"io"
	//"mime"
	"mime/multipart"	
)

func ProcessImage(image *multipart.FileHeader)(string,error){
	//open the file
   file,err:=image.Open()
   if err!=nil{
	  return "",err
   }
   defer file.Close()
  
   //Read the file 
   imageBytes,err:= io.ReadAll(file)
   if err!=nil{
	  return "",err
   }
   
   //Convert to base64
   base64Str:= base64.StdEncoding.EncodeToString(imageBytes)
   
   //Get the mime type example image/jpeg
   mimeType:=image.Header.Get("Content-Type")

   //final string given example data:image/jpeg;base64,/9j/4AAQSkZJRgABAQ...
   dataURI:="data:"+mimeType+";base64,"+base64Str

   return dataURI,nil
}