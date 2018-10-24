from server import app

from server.controllers.home_controller import HomeController
from server.controllers.login_controller import LoginController, LogoutController
from server.controllers.register_controller import RegisterController
from server.controllers.schedule_controller import ScheduleController
from server.controllers.course.search_controller import CourseSearchController, CourseEnrollingController
from server.controllers.course.enrolled_controller import CoursesEnrolledController

# Home page
home_view = HomeController.as_view('index')
app.add_url_rule('/', view_func=home_view)
app.add_url_rule('/home', view_func=home_view)

# Login and logout
app.add_url_rule('/login', view_func=LoginController.as_view('login'))
app.add_url_rule('/logout', view_func=LogoutController.as_view('logout'))

# Register
app.add_url_rule('/register', view_func=RegisterController.as_view('register'))

# Schedule page
app.add_url_rule('/schedule', view_func=ScheduleController.as_view('schedule'))

# Enrolled courses page
enrolled_courses_view = CoursesEnrolledController.as_view('courses_enrolled')
app.add_url_rule('/courses', view_func=enrolled_courses_view)
app.add_url_rule('/courses/enrolled', view_func=enrolled_courses_view)

# Enrolling in courses
app.add_url_rule('/courses/enroll', view_func=CourseEnrollingController.as_view('course_enrolling_hidden'))

# Course search page
app.add_url_rule('/courses/search', view_func=CourseSearchController.as_view('courses_search'))
