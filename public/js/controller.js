var askeecsApp = angular.module('askeecs', []);
 
askeecsApp.controller('QuestionListCtrl', function ($scope) {
  $scope.questions = [
    {'title': 'Nexus S', 'author': 'tperson'}
  ];
});
