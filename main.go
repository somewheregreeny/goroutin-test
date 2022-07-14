package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

type Root struct {
	XMLName    xml.Name   `xml:"Root"`
	Parameters Parameters `xml:"Parameters"`
	Dataset    Dataset    `xml:"Dataset"`
}

type Parameters struct {
	Parameter string `xml:"Parameter"`
}

type Dataset struct {
	ColumnInfo ColumnInfo `xml:"ColumnInfo"`
	Rows       Rows       `xml:"Rows"`
}

type ColumnInfo struct {
	Column []string `xml:"Column"`
}

type Rows struct {
	Row Row `xml:"Row"`
}

type Row struct {
	Col []Col `xml:"Col"`
}

type Col struct {
	Id   string `xml:"id,attr"`
	Data string `xml:",chardata"`
}

func send(yy, tmGbn, livstuNo, outStayGbn, outStayFrDt, outStayToDt, outStayAplyDt string,
	wg *sync.WaitGroup, client *http.Client) {
	sendStayOutXML := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
    <Root xmlns="http://www.nexacroplatform.com/platform/dataset">
        <Parameters>
            <Parameter id="_ga">GA1.3.1065330987.1626699518</Parameter>
            <Parameter id="requestTimeStr">1626795331154</Parameter>
        </Parameters>
        <Dataset id="DS_DORM120">
            <ColumnInfo>
                <Column id="chk" type="STRING" size="256"  />
                <Column id="yy" type="STRING" size="256"  />
                <Column id="tmGbn" type="STRING" size="256"  />
                <Column id="livstuNo" type="STRING" size="256"  />
                <Column id="outStaySeq" type="STRING" size="256"  />
                <Column id="outStayGbn" type="STRING" size="256"  />
                <Column id="outStayFrDt" type="STRING" size="256"  />
                <Column id="outStayToDt" type="STRING" size="256"  />
                <Column id="outStayStGbn" type="STRING" size="256"  />
                <Column id="outStayStNm" type="STRING" size="256"  />
                <Column id="outStayAplyDt" type="STRING" size="256"  />
                <Column id="outStayReplyCtnt" type="STRING" size="256"  />
                <Column id="schregNo" type="STRING" size="256"  />
                <Column id="hldyYn" type="STRING" size="256"  />
                <Column id="resprHldyYn" type="STRING" size="256"  />
            </ColumnInfo>
            <Rows>
                <Row type="insert">
                    <Col id="yy">%s</Col>
                    <Col id="tmGbn">%s</Col>
                    <Col id="livstuNo">%s</Col>         
                    <Col id="outStayGbn">%s</Col>       
                    <Col id="outStayFrDt">%s</Col> 
                    <Col id="outStayToDt">%s</Col> 
                    <Col id="outStayStGbn">1</Col>     
                    <Col id="outStayStNm">미승인</Col>
                    <Col id="outStayAplyDt">%s</Col>
                </Row>
            </Rows>
        </Dataset>
    </Root>`,
		yy,
		tmGbn,
		livstuNo,
		outStayGbn,
		outStayFrDt,
		outStayToDt,
		outStayAplyDt,
	))

	req, err := http.NewRequest("POST", "https://dream.tukorea.ac.kr/aff/dorm/DormCtr/saveOutAplyList.do?menuId=MPB0022&pgmId=PPB0021", bytes.NewBuffer(sendStayOutXML))
	if err != nil {
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	wg.Done()
}

func main() {
	loginInfo := url.Values{
		"internalId": {""},
		"internalPw": {""},
		"gubun":      {"inter"},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "https://ksc.tukorea.ac.kr/sso/login_proc.jsp?returnurl=null", bytes.NewBufferString(loginInfo.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{
		Jar: jar,
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	req, err = http.NewRequest("GET", "https://dream.tukorea.ac.kr/com/SsoCtr/initPageWork.do?loginGbn=sso&loginPersNo=", nil)

	if err != nil {
		panic(err)
	}

	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	findUserNmXML := []byte(
		`<?xml version="1.0" encoding="UTF-8"?>
			<Root xmlns="http://www.nexacroplatform.com/platform/dataset">
				<Parameters>
					<Parameter id="columnList">persNo|userNm</Parameter>
					<Parameter id="requestTimeStr">1627027228674</Parameter>
				</Parameters>
			</Root>`)

	req, err = http.NewRequest("POST", "https://dream.tukorea.ac.kr/com/SsoCtr/findMyGLIOList.do?menuId=MPB0022&pgmId=PPB0021", bytes.NewBuffer(findUserNmXML))
	if err != nil {
		panic(err)
	}

	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	var studentInfo Root
	xml.Unmarshal(body, &studentInfo)

	findYYtmgbnXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Root xmlns="http://www.nexacroplatform.com/platform/dataset">
	<Parameters>
		<Parameter id="_ga">GA1.3.1065330987.1626699518</Parameter>
		<Parameter id="requestTimeStr">1626878279321</Parameter>
	</Parameters>
	<Dataset id="DS_COND">
		<ColumnInfo>
			<Column id="mvinTermYn" type="STRING" size="256"  />
		</ColumnInfo>
		<Rows>
			<Row>
				<Col id="mvinTermYn">1</Col>
			</Row>
		</Rows>
	</Dataset>
</Root>`)

	req, err = http.NewRequest("POST", "https://dream.tukorea.ac.kr/aff/dorm/DormCtr/findYyTmGbnList.do?menuId=MPB0022&pgmId=PPB0021", bytes.NewBuffer(findYYtmgbnXML))
	if err != nil {
		panic(err)
	}

	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ = ioutil.ReadAll(res.Body)
	var yytmGbnInfo Root
	xml.Unmarshal(body, &yytmGbnInfo)

	findLiveStuNoXML := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
    <Root xmlns="http://www.nexacroplatform.com/platform/dataset">
        <Parameters>
            <Parameter id="_ga">GA1.3.1065330987.1626699518</Parameter>
            <Parameter id="requestTimeStr">1626877490927</Parameter>
        </Parameters>
        <Dataset id="DS_COND">
            <ColumnInfo>
                <Column id="yy" type="STRING" size="256"  />
                <Column id="tmGbn" type="STRING" size="256"  />
                <Column id="schregNo" type="STRING" size="256"  />
                <Column id="stdKorNm" type="STRING" size="256"  />
                <Column id="outStayStGbn" type="STRING" size="256"  />
            </ColumnInfo>
            <Rows>
                <Row type="update">
                    <Col id="yy">%s</Col>
                    <Col id="tmGbn">%s</Col>
                    <Col id="schregNo">%s</Col>
                    <Col id="stdKorNm">%s</Col>
                    <OrgRow>
                    </OrgRow>
                </Row>
            </Rows>
        </Dataset>
    </Root>`,
		yytmGbnInfo.Dataset.Rows.Row.Col[0].Data,
		yytmGbnInfo.Dataset.Rows.Row.Col[1].Data,
		studentInfo.Dataset.Rows.Row.Col[1].Data,
		studentInfo.Dataset.Rows.Row.Col[0].Data))

	req, err = http.NewRequest("POST", "https://dream.tukorea.ac.kr/aff/dorm/DormCtr/findMdstrmLeaveAplyList.do?menuId=MPB0022&pgmId=PPB0021", bytes.NewBuffer(findLiveStuNoXML))
	if err != nil {
		panic(err)
	}

	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ = ioutil.ReadAll(res.Body)
	var liveStuNo Root
	xml.Unmarshal(body, &liveStuNo)

	dateList := []string{"20220715", "20220716", "20220717", "20220718", "20220719", "20220720", "20220721", "20220722",
		"20220723", "20220724", "20220725", "20220726", "20220727", "20220728", "20220729", "20220730", "20220731",
		"20220801", "20220802", "20220803", "20220804", "20220805", "20220806", "20220807", "20220808", "20220809",
		"20220810", "20220811"}
	isWeekend := []int{0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0}
	outStayAplyDt := "20220715"

	var wg sync.WaitGroup
	var outStayGbn string

	wg.Add(len(dateList))

	for i := 0; i < len(dateList); i++ {
		if isWeekend[i] == 0 {
			outStayGbn = "07"
		} else {
			outStayGbn = "04"
		}

		go send(
			yytmGbnInfo.Dataset.Rows.Row.Col[0].Data,
			yytmGbnInfo.Dataset.Rows.Row.Col[1].Data,
			liveStuNo.Dataset.Rows.Row.Col[12].Data,
			outStayGbn,
			dateList[i],
			dateList[i],
			outStayAplyDt,
			&wg,
			client,
		)
	}

	wg.Wait()
	fmt.Println("Done")
}
