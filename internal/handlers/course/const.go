package course

var queryStr = `
        SELECT 
            MAX(ci.id) AS id,
            ci.course_num,
        	ti.classroom,
            ti.building,
            ANY_VALUE(ci.average_rating),
            ANY_VALUE(ci.review_count),
            ANY_VALUE(ci.credit),
            ANY_VALUE(ci.course_type) AS course_type,
            ANY_VALUE(ci.faculty) AS faculty,
            ANY_VALUE(ci.course_name) AS course_name,
            ANY_VALUE(ci.teacher) AS teacher,
            ANY_VALUE(ci.teacher_title) AS teacher_title,
            ANY_VALUE(ti.week_and_time) AS week_and_time,
            ANY_VALUE(ti.day_of_week) AS day_of_week
        FROM time_infos ti 
        JOIN course_infos ci ON ci.id = ti.course_info_id
        WHERE ti.day_of_week = ? 
          AND ti.area = ? 
          AND ti.week_and_time & ? = ?
        GROUP BY 
            ti.building, 
            ti.classroom,
            ci.course_num
    `
