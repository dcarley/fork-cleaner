package main_test

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/google/go-github/github"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

func repoFixture(name string, fork bool, user *github.User) *github.Repository {
	fullname := *user.Login + "/" + name

	return &github.Repository{
		Name:     github.String(name),
		FullName: github.String(fullname),
		HTMLURL:  github.String("https://github.com/" + fullname),
		Fork:     github.Bool(fork),
		Owner:    user,
	}
}

var _ = Describe("Main", func() {
	const username = "myusername"

	var (
		token   string
		stdin   *bytes.Buffer
		server  *ghttp.Server
		session *gexec.Session
		user    *github.User
		repos   []*github.Repository
	)

	BeforeEach(func() {
		stdin = new(bytes.Buffer)

		user = &github.User{
			Login: github.String(username),
		}

		status := http.StatusOK
		server = ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/user"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, user),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/users/"+username+"/repos"),
				ghttp.RespondWithJSONEncodedPtr(&status, &repos),
			),
		)
	})

	JustBeforeEach(func() {
		command := exec.Command(cmdPath, fmt.Sprintf("-baseURL=%s/", server.URL()))
		command.Env = append(command.Env, fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", token))
		command.Stdin = stdin

		var err error
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("token environment variable empty", func() {
		BeforeEach(func() {
			token = ""
		})

		It("should return error", func() {
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say(`Must set environment variable: GITHUB_ACCESS_TOKEN`))
		})
	})

	Describe("token environment variable set", func() {
		BeforeEach(func() {
			token = "mysekret"
		})

		Describe("two repos are not a fork", func() {
			BeforeEach(func() {
				repos = []*github.Repository{
					repoFixture("isfork1", true, user),
					repoFixture("notfork1", false, user),
					repoFixture("notfork2", false, user),
					repoFixture("isfork2", true, user),
				}

				stdin.WriteString("\n")
				stdin.WriteString("\n")
				stdin.WriteString("\n")
				stdin.WriteString("\n")
			})

			It("should only prompt for forks", func() {
				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say(`isfork1`))
				Eventually(session.Out).ShouldNot(gbytes.Say(`notfork1`))
				Eventually(session.Out).ShouldNot(gbytes.Say(`notfork2`))
				Eventually(session.Out).Should(gbytes.Say(`isfork2`))
			})
		})

		Describe("two repos requested to be deleted", func() {
			BeforeEach(func() {
				repos = []*github.Repository{
					repoFixture("fork1", true, user),
					repoFixture("fork2", true, user),
					repoFixture("fork3", true, user),
					repoFixture("fork4", true, user),
				}

				stdin.WriteString("fork1\n")
				stdin.WriteString("\n")
				stdin.WriteString("fork3\n")
				stdin.WriteString("\n")

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/repos/"+username+"/fork1"),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/repos/"+username+"/fork3"),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
				)
			})

			It("should delete two repos", func() {
				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say(`fork1 has been deleted`))
				Eventually(session.Out).ShouldNot(gbytes.Say(`fork2 has been deleted`))
				Eventually(session.Out).Should(gbytes.Say(`fork3 has been deleted`))
				Eventually(session.Out).ShouldNot(gbytes.Say(`fork4 has been deleted`))
			})
		})
	})
})
