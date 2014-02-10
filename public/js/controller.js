var askeecsControllers = angular.module('askeecsControllers', []);

askeecsControllers.controller('QuestionListCtrl', ['$scope', '$http',
	function ($scope, $http) {
		$http.get('data/questions.json').success(function(data) {
			$scope.questions = data;
		});
	}
]);

askeecsControllers.controller('RegisterCtrl', ['$scope', '$http',
	function ($scope, $http) {
		$scope.data = {}
		$scope.processForm = function () {
			console.log("GO!");
			if($scope.data.Password != $scope.data.cpassword) {
				console.log("Missed matched password");
				return;
			}

			delete $scope.data.cpassword;
			$scope.data.Username += '@email.wsu.edu'
			console.log($scope.data);
			$http({
				method: 'POST',
				url: '/register',
				data: $scope.data
			}).success(function(data) {

			});
			
		}
	}
]);

askeecsControllers.controller('QuestionAskCtrl', ['$scope', '$http', '$window', '$sce',
	function ($scope, $http, $window, $sce) {
		$scope.markdown="";
		$scope.title="";
		$scope.tags="";
		$scope.md2Html = function() {
			$scope.html = $window.marked($scope.markdown);
			$scope.htmlSafe = $sce.trustAsHtml($scope.html);
		}

		$scope.processForm = function () {
			console.log($scope.markdown);
			console.log($scope.tags);
			console.log($scope.title);
			delete $scope.errorMarkdown;
			delete $scope.errorTitle;
			delete $scope.errorTags;

			var err = false;

			if ($scope.markdown.length < 120)
			{
				$scope.errorMarkdown = "Your question must be 120 characters or more."
				err = true;
			}

			if ($scope.title.length == 0)
			{
				$scope.errorTitle = "You must enter a title."
				err = true;
			}

			if ($scope.tags.length == 0)
			{
				$scope.errorTags = "You must have at least one tag."
				err = true;
			}




			if (err) {
				return;
			}

			$http({
				method: 'POST',
				url: '/question',
				data: $.param({title:$scope.title, body: $scope.markdown, tags: $scope.tags.split(' ')})
			}).success(function(data) {
				if (!data.success) {
					
				}
				else
				{

				}
			});
		}

	}
]);

askeecsControllers.controller('QuestionDetailCtrl', ['$scope', '$routeParams', '$http',
	function ($scope, $routeParams, $http) {
		$http.get('data/questions.json').success(function(data) {
			$scope.question = data[$routeParams.questionId];
		});
	}
]);
