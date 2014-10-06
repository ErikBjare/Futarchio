app = angular.module('FutarchioApp', ["ngResource", "ngRoute", "ngCookies"]);

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
   .when('/predictions', {
    templateUrl: 'predictions.html',
    controller: 'PredictionsController',
  })
   .when('/profile/:username', {
    templateUrl: 'profile.html',
    controller: 'ProfileController',
  })
   .when('/profile', {
    templateUrl: 'profile.html',
    controller: 'ProfileController',
  })
   .when('/admin', {
    templateUrl: 'admin.html',
    controller: 'AdminController',
  })
   .when('/logout', {
    templateUrl: 'logout.html',
    controller: 'LogoutController',
  })
   .when('/login', {
    templateUrl: 'login.html',
    controller: 'LoginController',
  });
});
