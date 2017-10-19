#!/bin/env groovy​

folder("triggers")

job("triggers/tectonic-installer-pr-trigger") {
  description('Tectonic Installer PR Trigger. Changes here will be reverted automatically.')

  concurrentBuild()

  logRotator(30, 100)
  label("worker && ec2")

  parameters {
    stringParam('ghprbPullId', '', 'PR number')
  }

  properties {
    githubProjectUrl('https://github.com/coreos/tectonic-installer')
  }

  wrappers {
    colorizeOutput()
    timestamps()
    buildInDocker {
      image('quay.io/coreos/tectonic-smoke-test-env:v5.6')
    }
    timeout {
        absolute(30)
    }
  }

  triggers {
    ghprbTrigger {
      gitHubAuthId("")
      adminlist("")
      orgslist("coreos\ncoreos-inc")
      whitelist("")
      cron("H/5 * * * *")
      triggerPhrase("ok to test")
      onlyTriggerPhrase(false)
      useGitHubHooks(true)
      permitAll(false)
      autoCloseFailedPullRequests(false)
      displayBuildErrorsOnDownstreamBuilds(false)
      commentFilePath("")
      skipBuildPhrase(".*\\[skip\\W+ci\\].*")
      blackListCommitAuthor("")
      allowMembersOfWhitelistedOrgsAsAdmin(true)
      msgSuccess("")
      msgFailure("")
      commitStatusContext("Jenkins-Tectonic-Installer")
      buildDescTemplate("#\$pullId: \$abbrTitle")
      blackListLabels("")
      whiteListLabels("")
      includedRegions("")
      excludedRegions("")
    }
  }

  steps {
    shell """#!/bin/bash -ex
      curl "https://api.github.com/repos/coreos/tectonic-installer/labels" > repoLabels
      repoLabels=\$(jq ".[] | .name" repoLabels)
      repoLabels=\$(echo \$repoLabels | tr -d "\\"" | tr [a-z] [A-Z] | tr - _)
      for label in \$repoLabels
      do
        echo \$label=false >> env_vars
      done


      curl "https://api.github.com/repos/coreos/tectonic-installer/issues/\${ghprbPullId}" > pr
      labels=\$(jq ".labels | .[] | .name" pr)
      labels=\$(echo \$labels | tr -d "\\"" | tr [a-z] [A-Z] | tr - _)
      for label in \$labels
      do
        echo \$label=true >> env_vars
      done
    """

    downstreamParameterized {
      trigger('tectonic-installer/PR-\${ghprbPullId}') {
        parameters {
          propertiesFile("env_vars", true)
        }
      }
    }
    shell "sleep 5"
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Installer PR Trigger")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      notifyRepeatedFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
    publishers {
        groovyPostBuild("""
import jenkins.model.Jenkins
import hudson.model.*
import org.jenkinsci.plugins.workflow.job.WorkflowRun
import org.jenkinsci.plugins.workflow.support.steps.StageStepExecution
import org.jenkinsci.plugins.workflow.job.WorkflowJob

//Get the PR Number
def thr = Thread.currentThread()
def currentBuild = thr?.executable
def resolver = currentBuild.buildVariableResolver
def PRNum = resolver.resolve("ghprbPullId")


// sleep a bit to wait jenkins refresh the jobs
sleep(5000);

def params = [ ];

// Get the PR Job
def job = Jenkins.instance.getItemByFullName("tectonic-installer/PR-" + PRNum)
manager.listener.logger.println("Jobs: " + job);
manager.listener.logger.println("PR Num: " + PRNum);
// If job is in the queue wait for that
manager.listener.logger.println("Job is in queue?: " + job.isInQueue());
while(job.isInQueue()) {
  manager.listener.logger.println("Job in the queue, waiting....");
  sleep(1000);
}
manager.listener.logger.println(job.builds);
for (prBuild in job.builds) {
manager.listener.logger.println("Job Num: " + prBuild.getNumber().toInteger());
manager.listener.logger.println("Job is building?: " + prBuild.isBuilding());
  if (prBuild.getNumber().toInteger() == 1 && prBuild.isBuilding()) {
    manager.listener.logger.println("Build 1 is running, will try to kill...");
    WorkflowRun run = (WorkflowRun) prBuild;
    //hard kill
    run.doKill();

    while(prBuild.isBuilding()) {
      manager.listener.logger.println("Trying to kill the job....");
      run.doKill();
      sleep(1000);
    }

    manager.listener.logger.println("Job Killed");
    //release pipeline concurrency locks
    StageStepExecution.exit(run);

    sleep(1000);
    def parameters = currentBuild?.actions.find{ it instanceof ParametersAction }?.parameters

    def cause = new Cause.UpstreamCause(currentBuild)
    def causeAction = new hudson.model.CauseAction(cause)

    def pr_trigger_job = Jenkins.instance.getItemByFullName("triggers/tectonic-installer-pr-trigger")
    def paramsAction = new ParametersAction(parameters)
    manager.listener.logger.println(parameters);

    hudson.model.Hudson.instance.queue.schedule(pr_trigger_job, 0, causeAction, paramsAction)
    break;
  }
}
manager.listener.logger.println("Done");
"""
    ,Behavior.MarkFailed)
    }
  }
}
