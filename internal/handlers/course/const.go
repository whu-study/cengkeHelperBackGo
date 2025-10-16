package course

var queryStr = `
        SELECT 
            MAX(ci.id) AS id,
            ci.course_num,
        	ti.classroom,
            ti.building,
            MAX(ci.course_type) AS course_type,
            MAX(ci.faculty) AS faculty,
            MAX(ci.course_name) AS course_name,
            MAX(ci.teacher) AS teacher,
            MAX(ci.teacher_title) AS teacher_title,
            MAX(ti.week_and_time) AS week_and_time,
            MAX(ti.day_of_week) AS day_of_week
        FROM time_infos ti 
        JOIN course_infos ci ON ci.id = ti.course_info_id
        WHERE ti.day_of_week = ? 
          AND ti.area = ? 
          AND (? = -1 OR (ti.week_and_time & (1 << (32 - ?))) != 0)
          AND (? = -1 OR (ti.week_and_time & (1 << (? - 1))) != 0)
        GROUP BY 
            ti.building, 
            ti.classroom,
            ci.course_num
    `
