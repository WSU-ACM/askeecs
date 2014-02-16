var askeecsControllers = angular.module('askeecsControllers', ['ngCookies']);

askeecsControllers.controller('QuestionListCtrl', ['$scope', '$http',
	function ($scope, $http) {
		$http.get('/q').success(function(data) {
			$scope.questions = data;
		});
	}
]);

askeecsControllers.controller('RegisterCtrl', ['$scope', '$http', '$location', 'AuthService',
	function ($scope, $http, $location, AuthService) {
		var credentials = { "Username": "", "Password": "", "CPassword": "" }

		$scope.credentials = credentials; 
		$scope.processForm = function () {

			// Make sure they have entered a password that matches
			if($scope.credentials.Password != $scope.credentials.CPassword) {
				console.log("Missed matched password");
				return;
			}

			// We don't need this to be passed along
			delete $scope.credentials.CPassword;

			// Register the user and redirect them to the login page
			AuthService.register($scope.credentials, function () {
				$location.path("/login");
			});

			// Make sure we wipe out the credentials
			$scope.credentials = credentials; 

		}
	}
]);

askeecsControllers.controller('LoginCtrl', ['$scope', '$http', '$location', 'AuthService',
	function ($scope, $http, $location, AuthService) {
		var credentials = { "Username": "", "Password": "", "Salt": "" }

		$scope.credentials = credentials
		$scope.processForm = function () {

			// Log the user in and direct them tot he home page
			AuthService.login($scope.credentials, function () {
				$location.path("/");
			});

			// Make sure we wipe out the credentials
			$scope.credentials = credentials
		}
	}
]);

askeecsControllers.controller('QuestionAskCtrl', ['$scope', '$http', '$window', '$sce', '$location',
	function ($scope, $http, $window, $sce, $location) {
		var question = {"markdown" : "", "title" : "", "tags" : ""}

		$scope.question = question;

		$scope.md2Html = function() {
			var src			= $scope.markdown || ""
			var html		= $window.marked(src);
			$scope.htmlSafe = $sce.trustAsHtml(html);
		}

		$scope.processForm = function () {

			// Remove any previous error statements
			$scope.error = {}


			// Default to a non error state
			var err = false;

			if ($scope.markdown.length < 50)
			{
				$scope.error.markdown = "Your question must be 50 characters or more."
				err = true;
			}

			if ($scope.title.length == 0)
			{
				$scope.error.title = "You must enter a title."
				err = true;
			}

			if ($scope.tags.length == 0)
			{
				$scope.error.tags = "You must have at least one tag."
				err = true;
			}

			if (err) {
				return;
			}

			$http({
				method: 'POST',
				url: '/q',
				data: {Title:$scope.title, Body: $scope.markdown, Tags: $scope.tags.split(' ')}
			}).success(function(data) {
				// TODO: this should be a JSON response
				$location.path("/questions/"+data);	
			});
			// TODO: Failure
		}

	}
]);

askeecsControllers.controller('QuestionDetailCtrl', ['$scope', '$routeParams', '$http', '$window', '$sce',
	function ($scope, $routeParams, $http, $window, $sce) {
		$scope.comment = { "Body" : "" }
		$scope.response = { "Body" : "" }

		$http.get('/q/' + $routeParams.questionId).success(function(data) {
			$scope.question = data;
			console.log(data)
		});

		$scope.voteUp = function () {
			$http({
				method: 'GET',
				url: '/q/' + $scope.question.ID + '/vote/up',
				data: {}
			}).success(function(data) {
				$scope.question.Upvotes = data.Upvotes
			});
		}

		$scope.voteDown = function () {
			$http({
				method: 'GET',
				url: '/q/' + $scope.question.ID + '/vote/down',
				data: {}
			}).success(function(data) {
				$scope.question.Downvotes = data.Downvotes
			});
		}

		$scope.markdown="";
		$scope.md2Html = function() {
			var src = $scope.response.Body || ""
			$scope.html = $window.marked(src);
			$scope.htmlSafe = $sce.trustAsHtml($scope.html);
		}

		// Can a comment have this own controller and it's own scope?
		$scope.processComment = function () {
			delete $scope.errorComment;

			var err = false;

			if ( $scope.comment.Body.length < 15 )
			{
				$scope.errorComment = "Your comment must be at least 15 characters"
				err = true;
			}

			if (err) return;

			$http({
				method: 'post',
				url: '/q/' + $scope.question.ID + '/comment/',
				data: $scope.comment
			}).success(function(data) {
				delete $scope.scomment;
				$scope.question.Comments.push(data);
			});
		}

		$scope.processForm = function () {
			console.log($scope.response.Body);
			delete $scope.errorMarkdown;

			var err = false;

			if ($scope.response.Body.length < 50)
			{
				$scope.errorMarkdown = "Your response must be 50 characters or more."
				err = true;
			}


			if (err) {
				return;
			}

			$http({
				method: 'post',
				url: '/q/' + $scope.question.ID + '/response/',
				data: $scope.response
			}).success(function(data) {
				$scope.question.Responses.push(data);
			});
		}
	}
]);
