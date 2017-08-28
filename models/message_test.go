package models

import (
    "testing"
)

func TestMessage(t *testing.T) {
    testMessage := Message{MessageBody:"test", Recipients:[]string{"12"}, Originator:"123", Encoding:NORMAL}

    if !testMessage.IsEncodingNormal() {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", "Normal", "otherwise")
    }

}

func TestMessage_GetMessagebodyLength(t *testing.T) {
    testMessage := Message{MessageBody:"test", Recipients:[]string{"12"}, Originator:"123", Encoding:NORMAL}
    actual := testMessage.GetMessagebodyLength()
    if actual != 4 {
        t.Errorf("Test failed, length expected: '%s', got:  '%s'", "4", actual)
    }
}

func TestMessage_GetSplitMessageWithOutBodyFromMessage(t *testing.T) {
    testMessage := Message{MessageBody:"test", Recipients:[]string{"12"}, Originator:"123", Encoding:NORMAL}
    sM := testMessage.GetSplitMessageWithOutBodyFromMessage()

    if sM.Originator != testMessage.Originator {
        t.Errorf("Test failed, Originator expected: '%s', got:  '%s'", testMessage.Originator, sM.Originator)
    }

    tM := Message{MessageBody:"test", Recipients:[]string{"12"}, Originator:"123", Encoding:UNICODE}
    sM2 := tM.GetSplitMessageWithOutBodyFromMessage()
    if sM2.DataCoding != "unicode" {
        t.Errorf("Test failed, Originator expected: '%s', got:  '%s'", "unicode", sM.DataCoding)
    }
}

func TestMessage_GetMessagebodyLength2(t *testing.T) {
    testMessage := Message{MessageBody:"test", Recipients:[]string{"12"}, Originator:"123", Encoding:UNICODE}
    actual := testMessage.GetMessagebodyLength()
    if actual != 4 {
        t.Errorf("Test failed, length expected: '%s', got:  '%s'", "4", actual)
    }
}
