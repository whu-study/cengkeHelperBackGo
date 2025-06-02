package course

var queryStr = `
        SELECT 
            ci.course_type,
            ci.faculty,
            ci.course_name,
            ci.teacher,
            ci.teacher_title,
            ti.classroom,
            ti.week_and_time,
            ti.building
        FROM time_infos ti 
        JOIN course_infos ci ON ci.id = ti.course_info_id
        WHERE ti.day_of_week = ? AND ti.area = ? 
        GROUP BY 
            ti.building, 
            ti.classroom, 
            ci.course_type, 
            ci.faculty, 
            ci.course_name, 
            ci.teacher, 
            ci.teacher_title,
            ti.week_and_time
    `
