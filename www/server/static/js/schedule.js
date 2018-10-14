$(function () {
    var $timeline = $("#timeline")
    var $timeline_text = $("#timeline-text");

    function hideTimeline() {
        $timeline.css('display', 'none')
        $timeline_text.css('display', 'none')
    }

    function showTimeline() {
        $timeline.css('display', 'inline')
        $timeline_text.css('display', 'inline')
    }

    function setTimelineY(y) {
        $timeline.attr('y1', y + "px")
        $timeline.attr('y2', y + "px")

        var time = (y - 50);
        time /= 60;

        var hour = Math.floor(time);
        var minutes = Math.floor((time - hour) * 60)
        hour += 8;

        var am_or_pm = "AM";

        if (hour > 12) {
            hour -= 12;
            am_or_pm = "PM";
        }

        if (minutes < 10) {
            minutes = "0" + minutes.toString();
        } else {
            minutes = minutes.toString();
        }

        $timeline_text.attr('y', (y - 2) + "px")
        $timeline_text.html(hour + ":" + minutes + " " + am_or_pm);
    }

    $("#schedule-svg").on("mousemove", function (e) {
        if (e.offsetY > 50) {
            showTimeline();
            setTimelineY(e.offsetY);
        } else {
            hideTimeline();
        }
    });

    $("#schedule-svg").on("mouseleave", function (e) {
        hideTimeline();
    });
});