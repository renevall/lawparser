(function () {
    'use strict';

    angular
        .module('app.home')
        .controller('Home', Home);


    function Home() {
        // var vm = this;
        // vm.customers = [];
        // vm.gotoCustomer = gotoCustomer;
        // vm.title = 'Home';

        // activate();

        // function activate() {
        //     return getCustomers().then(function () {
        //         logger.info('Activated Home View');
        //     });
        // }

        // function getCustomers() {
        //     return dataservice.getCustomers().then(function (data) {
        //         vm.customers = data;
        //         return vm.customers;
        //     });
        // }

        // function gotoCustomer(c) {
        //     $state.go('customer.detail', {id: c.id});
        // }
    }
})();