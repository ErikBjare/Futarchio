# -*- coding: utf-8 -*-
from selenium import selenium
import unittest, time, re

class signup(unittest.TestCase):
    def setUp(self):
        self.verificationErrors = []
        self.selenium = selenium("localhost", 4444, "*chrome", "http://localhost:8080")
        self.selenium.start()
    
    def test_signup(self):
        sel = self.selenium
        sel.open("/")
        sel.click("link=Sign up")
        time.sleep(0.5)
        sel.type("//input[@type='username']", "tester")
        sel.type("//input[@type='password']", "testing")
        sel.type("//input[@type='text']", "tester@example.com")
        sel.click("//button[@type='submit']")
        time.sleep(0.5)
        self.assertEqual("http://localhost:8080/signup/success", sel.get_location())
        sel.open("/login")
        sel.type("name=username", "tester")
        sel.type("name=password", "testing")
        sel.click("//button[@type='submit']")
        sel.wait_for_page_to_load("30000")
        sel.click("link=My Account")
        sel.click("link=Log out")
        sel.wait_for_frame_to_load("", "")
        sel.click("link=Sign up")
        sel.open("/")
    
    def tearDown(self):
        self.selenium.stop()
        self.assertEqual([], self.verificationErrors)

if __name__ == "__main__":
    unittest.main()
