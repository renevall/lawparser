import { Component, View } from "angular2/core";
import { bootstrap } from "angular2/platform/browser";

@Component({
    selector: "my-app",
    templateUrl: "app/app.html",
    directives: []
})

class App {

  filesToUpload: Array<File>;

  constructor() {
    this.filesToUpload = [];
  }

  upload(){
    this.makeFileRequest("http://localhost:8080/upload",[],
    this.filesToUpload).then((result) => {
      console.log(result);
    }, (error) => {
      console.log(error);
    });
  }

  fileChangeEvent(fileInput: any){
    this.filesToUpload = <Array<File>> fileInput.target.files;
  }

  makeFileRequest(url: string, params: Array<string>, files: Array<File>){
    console.log("making request");
    console.log(files);
    return new Promise((resolve, reject) => {
      var formData: any = new FormData();
      var xhr = new XMLHttpRequest();
      for(var i = 0; i < files.length; i++){
        formData.append("uploads[]", files[i], files[i].name);
      }
      console.log(formData);
      xhr.onreadystatechange = function(){
        if(xhr.readyState == 4){
          if(xhr.readyState == 200){
            resolve(JSON.parse(xhr.response));
          }else{
            reject(xhr.response);
          }
        }
      }
      xhr.open("POST", url, true);
      xhr.send(formData);
    });
  }

}

bootstrap(App, []);
