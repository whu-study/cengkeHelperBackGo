package course

var queryStr = `
        SELECT 
            ci.id,
            ci.course_type,
            ci.faculty,
            ci.course_name,
            ci.teacher,
            ci.teacher_title,
            ti.classroom,
            ti.week_and_time,
            ti.building,
            ti.day_of_week
        FROM time_infos ti 
        JOIN course_infos ci ON ci.id = ti.course_info_id
        WHERE ti.day_of_week = ? AND ti.area = ? 
        GROUP BY 
            ci.id,
            ti.building, 
            ti.classroom, 
            ci.course_type, 
            ci.faculty, 
            ci.course_name, 
            ci.teacher, 
            ci.teacher_title,
            ti.week_and_time,
            ti.day_of_week
    `
