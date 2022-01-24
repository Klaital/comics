package comics

func (suite *ComicTestSuite) TestMoveOrdinalsNoCollision() {
	db, err := suite.Cfg.ConnectDB()
	suite.NoError(err, "expected to connect to a test DB")
	suite.NotNil(db, "expected ConnectDB go return a connection")

	// Existing fixtures include ordinals 1-3.

	// Inserting at #4 should leave the existing entries untouched.
	err = MoveComicsDownIfCollision(4, 1, db)
	suite.NoError(err, "expected no error")

}
