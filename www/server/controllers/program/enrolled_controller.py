from server import controllers
from server.facades import program_facade
from server.models.user_model import UserModel
from server.models.program_model import ProgramModel
from server.models.utils import facade_result_codes as FRC
import flask


class ProgramsEnrolledController(controllers.BaseController):
    decorators = [controllers.login_required]

    def get(self):
        user = UserModel.create(**flask.session['user'])

        context = {
            "selected_tab": "programs",
            "programs": [],
        }

        (status, program_list) = program_facade.get_enrolled_programs(user.token)

        if status == FRC.NOT_AUTHENTICATED:
            return self.handle_not_authorized()

        if status != FRC.SUCCESS:
            flask.flash(f"Error loading your currently enrolled programs.", "danger")
        else:
            program_list.append(ProgramModel.create(
                program_id=1,
                catalog_year=2018,
                major=1,
                program="Computer Science (B.S.)",
            ))

            context['programs'] = program_list

        return flask.render_template('programs/enrolled.html', **context)
