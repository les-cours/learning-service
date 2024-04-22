
CREATE TABLE classrooms (
    classroom_id Text PRIMARY KEY,
    teacher_id varchar(50),
    subject_id varchar(50),
);

CREATE TABLE videos (
    video_id varchar(50) PRIMARY KEY,
    title varchar(255),
);

CREATE TABLE comments (
    comment_id varchar(50) PRIMARY KEY,
    created_by varchar(50) ,
    is_owner BOOLEAN, -- if the owner of video comment
    video_id varchar(50) ,
    comment varchar(255),
    created_at timestamp
    replyto varchar(50)
);

CREATE TABLE pdfs (
    pdf_id varchar(50) PRIMARY KEY,
    title varchar(255),
);

CREATE TABLE classrooms_document (
    classroom_id varchar(50) PRIMARY KEY,
    document_id varchar(50),  --video_id , pdf_id
    document_type varchar(6), --(pdf/video)
);

-- subscription
CREATE TABLE classrooms_students (
    classroom_id varchar(50) PRIMARY KEY;
    student_id varchar(50),
);
