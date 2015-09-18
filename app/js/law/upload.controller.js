(function(){
	'use strict';

	angular
		.module('app.law')
		.controller('UploadController', UploadController);

	UploadController.$inject = ['Upload', '$timeout'];
	function UploadController(Upload, $timeout){
		var vm = this;

		activate();

		function activate(){

		}

		vm.uploadFile = function(file){
			console.log(file);
			console.log(vm.lawname)
			// console.log("entro");
			file.upload = Upload.upload({
				url: 'http://localhost:8080/api/law/uploadfile',
				method: 'POST',
				headers: {
					'my-header': 'my-header-value',
				},
				fields: {filename: vm.lawname},
				file: file,
				fileFormDataName: 'myFile'

			});

			file.upload.then(function(response){
				console.log(response);
				$timeout(function () {
        		file.result = response.data;
      			});
			}, function(response){
				if (response.status > 0)
				vm.errorMsg = response.status + ': ' +
					response.data;
			});

			file.upload.progress(function (evt) {
      		// Math.min is to fix IE which reports 200% sometimes
      			file.progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
    		});
		}
	}



})();