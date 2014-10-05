app = angular.module('WeArePeopleApp', ["ngResource", "ngRoute"]);

app.config(function($routeProvider, $locationProvider) {
  $routeProvider
   .when('/', {
    templateUrl: 'home.html',
    controller: 'HomeController',
  })
   .when('/polls', {
    templateUrl: 'polls.html',
    controller: 'PollsController',
  })
   .when('/profile/:username', {
    templateUrl: 'profile.html',
    controller: 'ProfileController',
  })
   .when('/admin', {
    templateUrl: 'admin.html',
    controller: 'AdminController',
  })
   .when('/login', {
    templateUrl: 'login.html',
    controller: 'LoginController',
  });
});
