app = angular.module('FutarchioApp', ["ngResource", "ngAnimate", "ngRoute", "ngCookies", "ngSanitize"]);

app.config(function($routeProvider, $locationProvider) {
  $routeProvider

  // Home
  .when('/', {
    templateUrl: '/static/home.html',
    controller: 'HomeController',
  })

  // Votes and polls
  .when('/vote', {
    templateUrl: '/static/polls.html',
    controller: 'PollsController',
  })
  .when('/poll/:key', {
    templateUrl: '/static/poll.html',
    controller: 'PollController',
  })

  // Predictions and statements
  .when('/predict', {
    templateUrl: '/static/predictions.html',
    controller: 'StatementsController',
  })
  .when('/statement/:key', {
    templateUrl: '/static/predictions.html',
    controller: 'StatementController',
  })
  .when('/prediction/:key', {
    templateUrl: '/static/predictions.html',
    controller: 'PredictionController',
  })

  // User profiles
  .when('/profile', {
   templateUrl: '/static/profile.html',
   controller: 'ProfileController',
  })
  .when('/profile/:username', {
   templateUrl: '/static/profile.html',
   controller: 'ProfileController',
  })

  // Admin
  .when('/admin', {
   templateUrl: '/static/admin.html',
   controller: 'AdminController',
  })

  // Login and logout
  .when('/logout', {
    templateUrl: '/static/logout.html',
    controller: 'LogoutController',
  })
  .when('/login', {
    templateUrl: '/static/login.html',
    controller: 'LoginController',
  })

  // Registration
  .when('/signup', {
    templateUrl: '/static/signup.html',
    controller: 'SignupController',
  })
  .when('/signup/success', {
    templateUrl: '/static/signup-success.html'
  });

  $locationProvider.html5Mode(true);
});

app.animation('.slide', function() {
    var NG_HIDE_CLASS = 'ng-hide';
    return {
        beforeAddClass: function(element, className, done) {
            if(className === NG_HIDE_CLASS) {
                element.slideUp(done);
            }
        },
        removeClass: function(element, className, done) {
            if(className === NG_HIDE_CLASS) {
                element.hide().slideDown(done);
            }
        }
    };
});
