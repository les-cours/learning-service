package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/les-cours/learning-service/api/learning"
	"github.com/les-cours/learning-service/api/users"
	"github.com/les-cours/learning-service/utils"
	"math/rand"
	"time"
)

func (s *Server) InitClassRooms(ctx context.Context, in *learning.IDRequest) (*learning.Notifications, error) {

	rows, err := s.DB.Query(`SELECT DISTINCT classroom_id FROM subscription WHERE student_id = $1`, in.Id)

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	var notifications = make([]*learning.Notification, 0)
	for _, id := range ids {

		res, err := s.GetCurrentSubscription(ctx, &learning.IDRequest{
			Id:     id,
			UserID: in.Id,
		})
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		var classroomTitle string
		err = s.DB.QueryRow(`SELECT arabic_title from classrooms where  classroom_id = $1`, id).Scan(&classroomTitle)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		if !res.Status {
			//add not subscription notification
			var notification = new(learning.Notification)
			notification.Id = utils.GenerateUUIDString()
			notification.Title = "unpaid"
			notification.Content = classroomTitle + " لم يتم دفعها، يرجى دفعها."

			notifications = append(notifications, notification)

			//remove it from chatrooms (change status)
			err = s.MongoDB.ChangeStudentStatus(ctx, id, in.Id, false)
			if err != nil {
				s.Logger.Error(err.Error())
			}

			continue

		}

		if res.RestDays < 4 {
			var notification = new(learning.Notification)
			notification.Id = utils.GenerateUUIDString()
			notification.Title = "Subscription Expiring Soon"
			notification.Content = fmt.Sprintf("الاشتراك في %s سينتهي خلال %d أيام. يرجى تجديده.", classroomTitle, res.RestDays)

			notifications = append(notifications, notification)
		}

	}

	return &learning.Notifications{
		Notifications: notifications,
	}, nil
}

func (s *Server) CreateClassRooms(ctx context.Context, in *learning.CreateClassRoomsRequest) (*learning.OperationStatus, error) {
	for _, subjectID := range in.SubjectIDs {
		_, err := s.CreateClassRoom(ctx, &learning.CreateClassRoomRequest{
			TeacherID: in.TeacherID,
			SubjectID: subjectID,
		})
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
	}
	return &learning.OperationStatus{Success: true}, nil
}

func (s *Server) CreateClassRoom(ctx context.Context, in *learning.CreateClassRoomRequest) (*learning.OperationStatus, error) {

	classRoomID := utils.GenerateUUIDString()
	/*
		GENERATE NAME :
	*/

	var subjectName string
	var classRoomArabicName string

	err := s.DB.QueryRow(`SELECT title ,title_ar FROM subjects 
                         WHERE subject_id = $1 
`, in.SubjectID).Scan(&subjectName, &classRoomArabicName)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrNotFound("subject")
	}

	var classRoomName = subjectName

	/*
		INSERT INTO CLASSROOM THE NEW ONE ...
	*/
	_, err = s.DB.Exec(`INSERT INTO 
    classrooms 
    (classroom_id, teacher_id, subject_id, arabic_title,title,price,image,badge) 
     VALUES ($1,$2,$3,$4,$5,$6,'https://firebasestorage.googleapis.com/v0/b/uploadingfile-90303.appspot.com/o/images%2F20240620_162629.png?alt=media&token=fb7f692a-a336-4329-8050-461851ee99fb','new')`, classRoomID, in.TeacherID, in.SubjectID, classRoomArabicName, classRoomName, 799)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	room := &learning.ClassRoom{
		ClassRoomID:  classRoomID,
		Title:        classRoomName,
		Price:        799,
		StudentCount: 0,
		Rating:       0,
		ArabicTitle:  classRoomArabicName,
	}
	err = s.createChatRoom(ctx, room, in.TeacherID)
	if err != nil {
		s.Logger.Error(err.Error())
		//return nil, err
	}
	return &learning.OperationStatus{Success: true}, nil
}

func (s *Server) createChatRoom(ctx context.Context, room *learning.ClassRoom, accountID string) error {

	user, err := s.Users.GetUserByID(ctx, &users.GetUserByIDRequest{
		AccountID: accountID,
		UserRole:  "teacher",
	})
	var teacher = &learning.User{
		Id:        user.Id,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
	}
	if err != nil {
		s.Logger.Error(err.Error())
		teacher = &learning.User{
			Id:        "deleted_account",
			Username:  "deleted_account",
			FirstName: "deleted_account",
			LastName:  "deleted_account",
			Avatar:    "deleted_account",
		}
	}
	err = s.MongoDB.CreateRoom(ctx, room, teacher)
	if err != nil {
		s.Logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *Server) GetClassRooms(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var subjectID = in.Id

	//teacher_id uuid NULL,
	//subject_id character varying(40) NULL,

	rows, err := s.DB.Query(`
SELECT classroom_id, title, image, price, c.badge, c.description, arabic_title,c.description_ar,t.teacher_id,t.firstname,t.lastname
	FROM classrooms as c
	INNER JOIN public.teachers t on t.teacher_id = c.teacher_id
WHERE c.deleted_at IS NULL AND c.subject_id = $1 ;
        `, subjectID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	var classRooms = new(learning.ClassRooms)
	var studentsCount int32
	for rows.Next() {

		var classRoomID, title, image, badge, description, arabicTitle, arabicDescription, teacherID, firstname, lastname string
		var price int32

		err = rows.Scan(&classRoomID, &title, &image, &price, &badge, &description, &arabicTitle, &arabicDescription, &teacherID, &firstname, &lastname)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		classRoom := &learning.ClassRoom{
			ClassRoomID:       classRoomID,
			Title:             title,
			Image:             image,
			Price:             price,
			Badge:             badge,
			Description:       description,
			ArabicTitle:       arabicTitle,
			ArabicDescription: arabicDescription,
			Teacher: &learning.Teacher{
				TeacherID: teacherID,
				Firstname: firstname,
				Lastname:  lastname,
			},
		}

		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount
		//get Rating  //To Do

		classRoom.Rating = 4.7
		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

// GetClassRoom acting  like promo
func (s *Server) GetClassRoom(ctx context.Context, in *learning.IDRequest) (*learning.ClassRoom, error) {

	var classRoomID = in.Id
	var studentsCount int32
	var title, image, badge, description, arabicTitle, arabicDescription, teacherID, firstname, lastname string
	var price int32
	var err error

	err = s.DB.QueryRow(`
SELECT  title, image, price, c.badge, c.description, arabic_title,c.description_ar,t.teacher_id,t.firstname,t.lastname
	FROM classrooms as c
	INNER JOIN public.teachers t on t.teacher_id = c.teacher_id
WHERE c.deleted_at IS NULL AND c.classroom_id = $1 ;
        `, classRoomID).Scan(&title, &image, &price, &badge, &description, &arabicTitle, &arabicDescription, &teacherID, &firstname, &lastname)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	classRoom := &learning.ClassRoom{
		ClassRoomID:       classRoomID,
		Title:             title,
		Image:             image,
		Price:             price,
		Badge:             badge,
		Description:       description,
		ArabicTitle:       arabicTitle,
		ArabicDescription: arabicDescription,
		Teacher: &learning.Teacher{
			TeacherID: teacherID,
			Firstname: firstname,
			Lastname:  lastname,
		},
	}

	/*
		Get chapters
	*/

	chapters, err := s.GetChaptersByClassRoom(ctx, &learning.IDRequest{
		Id:     classRoomID,
		UserID: in.UserID,
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	classRoom.Chapters = chapters
	//get students count for this
	err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	classRoom.StudentCount = studentsCount
	//get Rating  //To Do

	classRoom.Rating = 4.7

	return classRoom, nil
}

func (s *Server) UpdateClassRoom(ctx context.Context, in *learning.UpdateClassRoomRequest) (*learning.ClassRoom, error) {

	err := userHasClassRoom(s.DB, in.TeacherID, in.ClassRoomID)
	if err != nil {
		return nil, err
	}
	_, err = s.DB.Exec(`UPDATE classrooms SET title = $2, image = $3, price = $4, arabic_title = $5, description = $6, description_ar = $7 WHERE classroom_id = $1;
`, in.ClassRoomID, in.Title, in.Image, in.Price, in.ArabicTitle, in.Description, in.ArabicDescription)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrNotFound("subject")
	}

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return s.GetClassRoom(ctx, &learning.IDRequest{Id: in.ClassRoomID})
}

func (s *Server) DeleteClassRoom(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {
	err := userHasClassRoom(s.DB, in.UserID, in.Id)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	var classRoomID = in.Id

	_, err = s.DB.Exec(`
UPDATE 
    classrooms 
SET deleted_at = $1 
WHERE classroom_id = $1;
        `, classRoomID)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &learning.OperationStatus{
		Success: true,
	}, nil
}

func (s *Server) DeleteClassRoomsByTeacher(ctx context.Context, in *learning.IDRequest) (*learning.OperationStatus, error) {

	_, err := s.DB.Exec(`
UPDATE 
    classrooms 
SET deleted_at = $1 
WHERE teacher_id = $1;
        `, in.Id)

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	return &learning.OperationStatus{
		Success: true,
	}, nil
}

func (s *Server) AddDocumentToClassroom(ctx context.Context, in *learning.AddDocumentToClassroomRequest) (*learning.OperationStatus, error) {

	return nil, nil
}

func (s *Server) GetClassRoomsByTeacher(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var teacherID = in.Id

	rows, err := s.DB.Query(`
SELECT classroom_id, title, image, price, badge, description, arabic_title
	FROM classrooms
WHERE deleted_at IS NULL AND teacher_id = $1;
        `, teacherID)
	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("classroom")
		}
		return nil, ErrInternal
	}
	classRooms := new(learning.ClassRooms)
	var studentsCount int32
	for rows.Next() {
		classRoom := &learning.ClassRoom{}
		err = rows.Scan(&classRoom.ClassRoomID, &classRoom.Title, &classRoom.Image, &classRoom.Price, &classRoom.Badge, &classRoom.Description, &classRoom.ArabicTitle)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount
		//get Rating  //To Do

		classRoom.Rating = 4.7

		chapters, err := s.GetChaptersByClassRoom(ctx, &learning.IDRequest{
			Id:     classRoom.ClassRoomID,
			UserID: teacherID,
		})
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		classRoom.Chapters = chapters

		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

func (s *Server) GetClassRoomsByStudent(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var subjectID = in.Id
	var studentID = in.UserID

	rows, err := s.DB.Query(`
SELECT classrooms.classroom_id, title, image, price, badge, description, arabic_title
	FROM classrooms
	INNER JOIN subscription s on classrooms.classroom_id = s.classroom_id
WHERE deleted_at IS NULL AND subject_id = $1 AND s.student_id = $2;
        `, subjectID, studentID)
	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("classroom")
		}
		return nil, ErrInternal
	}

	classRooms := &learning.ClassRooms{}
	var studentsCount int32
	for rows.Next() {
		classRoom := &learning.ClassRoom{}
		err = rows.Scan(&classRoom.ClassRoomID, &classRoom.Title, &classRoom.Image, &classRoom.Price, &classRoom.Badge)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount
		//get Rating  //To Do

		classRoom.Rating = 4.7
		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

//func (s *Server) DeleteClassRoom(ctx context.Context, in *learning.GetClassRoomsByStudentRequest) (*learning.ClassRooms, error) {
//
// 1 delete classrooms
//  2 delete videos and pdf from storage ...
//
//	return nil, nil
//}

/*
STUDENTS
*/

func (s *Server) GetMyClassRooms(ctx context.Context, in *learning.IDRequest) (*learning.ClassRooms, error) {

	var subjectID = in.Id
	var studentID = in.UserID

	currentTime := time.Now()

	rows, err := s.DB.Query(`
SELECT DISTINCT classrooms.classroom_id, title, COALESCE(image, 'default.jpg') AS image,price,COALESCE(badge, '') AS badge, description, arabic_title
	FROM classrooms
	INNER JOIN subscription s on classrooms.classroom_id = s.classroom_id
WHERE s.student_id = $1 AND deleted_at IS NULL AND classrooms.subject_id = $2 
AND s.paid_at between $3 AND $4;
        `, studentID, subjectID, currentTime.
		AddDate(0, -1, 0), currentTime)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	classRooms := &learning.ClassRooms{}
	var studentsCount int32

	var statistics = new(learning.Statistic)

	for rows.Next() {
		classRoom := &learning.ClassRoom{}
		err = rows.Scan(&classRoom.ClassRoomID, &classRoom.Title, &classRoom.Image, &classRoom.Price, &classRoom.Badge, &classRoom.Description, &classRoom.ArabicTitle)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		//get students count for this
		err = s.DB.QueryRow(`SELECT count(1) FROM subscription WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&studentsCount)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		classRoom.StudentCount = studentsCount

		_ = s.DB.QueryRow(`SELECT count(1) FROM chapters WHERE classroom_id = $1`, classRoom.ClassRoomID).Scan(&statistics.Chapters)

		statistics.Lessons = int32(rand.Intn(21))
		statistics.Videos = int32(rand.Intn(21))
		statistics.Pdfs = int32(rand.Intn(21))
		statistics.Students = studentsCount

		classRoom.Statistics = statistics
		classRoom.Rating = rand.Float32() * 5
		classRooms.Classrooms = append(classRooms.Classrooms, classRoom)
	}

	return classRooms, nil
}

func (s *Server) CanAccessToClassRoom(studentID, classroomID string) bool {

	var currentTime = time.Now()

	var exist bool
	err := s.DB.QueryRow(`
SELECT true
FROM subscription
WHERE paid_at  between $1 and $2
  AND student_id = $3
  AND classroom_id = $4 LIMIT 1`, currentTime.AddDate(0, -1, 0), currentTime, studentID, classroomID).Scan(&exist)

	if err != nil {
		s.Logger.Error(err.Error())
		return false
	}

	return true

}

func userHasClassRoom(db *sql.DB, userID, classRoomID string) error {
	var has bool
	err := db.QueryRow(`SELECT EXISTS (
	    SELECT 1
	    FROM classrooms
	    WHERE teacher_id = $1
	    AND classroom_id = $2
	);`, userID, classRoomID).Scan(&has)

	if err != nil {
		return ErrInternal
	}

	if has {
		return nil
	}

	return ErrPermission
}
