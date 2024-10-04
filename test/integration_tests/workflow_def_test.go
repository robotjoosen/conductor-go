//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
//  the License. You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
//  specific language governing permissions and limitations under the License.

package integration_tests

import (
	"context"
	"os"
	"testing"

	"github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/conductor-sdk/conductor-go/sdk/workflow"
	"github.com/conductor-sdk/conductor-go/test/testdata"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.ErrorLevel)
}

func TestHttpTask(t *testing.T) {
	httpTaskWorkflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_HTTP").
		OwnerEmail("test@orkes.io").
		Version(1).
		WorkflowStatusListenerEnabled(true).
		Add(testdata.TestHttpTask)
	err := testdata.ValidateWorkflow(httpTaskWorkflow, testdata.WorkflowValidationTimeout, model.CompletedWorkflow)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.ValidateWorkflowBulk(httpTaskWorkflow, testdata.WorkflowValidationTimeout, testdata.WorkflowBulkQty)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(httpTaskWorkflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func SimpleTask(t *testing.T) {
	err := testdata.ValidateTaskRegistration(*testdata.TestSimpleTask.ToTaskDef())
	if err != nil {
		t.Fatal(err)
	}
	simpleTaskWorkflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_SIMPLE").
		Version(1).
		Add(testdata.TestSimpleTask)
	err = testdata.TaskRunner.StartWorker(
		testdata.TestSimpleTask.ReferenceName(),
		testdata.SimpleWorker,
		testdata.WorkerQty,
		testdata.WorkerPollInterval,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.ValidateWorkflow(simpleTaskWorkflow, testdata.WorkflowValidationTimeout, model.CompletedWorkflow)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.ValidateWorkflowBulk(simpleTaskWorkflow, testdata.WorkflowValidationTimeout, testdata.WorkflowBulkQty)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.TaskRunner.DecreaseBatchSize(
		testdata.TestSimpleTask.ReferenceName(),
		testdata.WorkerQty,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(simpleTaskWorkflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func SimpleTaskWithoutRetryCount(t *testing.T) {
	taskToRegister := testdata.TestSimpleTask.ToTaskDef()
	taskToRegister.RetryCount = 0
	err := testdata.ValidateTaskRegistration(*taskToRegister)
	if err != nil {
		t.Fatal(err)
	}
	simpleTaskWorkflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_SIMPLE").
		Version(1).
		Add(testdata.TestSimpleTask)
	err = testdata.TaskRunner.StartWorker(
		testdata.TestSimpleTask.ReferenceName(),
		testdata.SimpleWorker,
		testdata.WorkerQty,
		testdata.WorkerPollInterval,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.ValidateWorkflow(simpleTaskWorkflow, testdata.WorkflowValidationTimeout, model.CompletedWorkflow)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.ValidateWorkflowBulk(simpleTaskWorkflow, testdata.WorkflowValidationTimeout, testdata.WorkflowBulkQty)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.TaskRunner.DecreaseBatchSize(
		testdata.TestSimpleTask.ReferenceName(),
		testdata.WorkerQty,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(simpleTaskWorkflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestInlineTask(t *testing.T) {
	inlineTaskWorkflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_INLINE_TASK").
		Version(1).
		Add(testdata.TestInlineTask)
	err := testdata.ValidateWorkflow(inlineTaskWorkflow, testdata.WorkflowValidationTimeout, model.CompletedWorkflow)
	if err != nil {
		t.Fatal(err)
	}
	err = testdata.ValidateWorkflowBulk(inlineTaskWorkflow, testdata.WorkflowValidationTimeout, testdata.WorkflowBulkQty)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(inlineTaskWorkflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestSqsEventTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_EVENT_SQS").
		Version(1).
		Add(testdata.TestSqsEventTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestConductorEventTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_EVENT_CONDUCTOR").
		Version(1).
		Add(testdata.TestConductorEventTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestNatsEventTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_EVENT_NATS").
		Version(1).
		Add(testdata.TestNatsEventTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestAmqpQueueEventTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_EVENT_AMQP_QUEUE").
		Version(1).
		Add(testdata.TestAmqpQueueEventTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestAmqpExchangeEventTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_EVENT_AMQP_EXCHANGE").
		Version(1).
		Add(testdata.TestAmqpExchangeEventTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestKafkaPublishTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_KAFKA_PUBLISH").
		Version(1).
		Add(testdata.TestKafkaPublishTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestDoWhileTask(t *testing.T) {

}

func TestTerminateTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_TERMINATE").
		Version(1).
		Add(testdata.TestTerminateTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestSwitchTask(t *testing.T) {
	workflow := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("TEST_GO_WORKFLOW_SWITCH").
		Version(1).
		Add(testdata.TestSwitchTask)
	err := testdata.ValidateWorkflowRegistration(workflow)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ValidateWorkflowDeletion(workflow)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func TestDynamicForkWorkflow(t *testing.T) {
	wf := workflow.NewConductorWorkflow(testdata.WorkflowExecutor).
		Name("dynamic_workflow_array_sub_workflow").
		Version(1).
		Add(createDynamicForkTask())
	err := wf.Register(true)
	if err != nil {
		t.Fatal()
	}

	err = testdata.ValidateWorkflowDeletion(wf)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func createDynamicForkTask() *workflow.DynamicForkTask {
	return workflow.NewDynamicForkTaskWithoutPrepareTask(
		"dynamic_workflow_array_sub_workflow",
	).Input(
		"forkTaskWorkflow", "extract_user",
	).Input(
		"forkTaskInputs", []map[string]interface{}{
			{
				"input": "value1",
			},
			{
				"sub_workflow_2_inputs": map[string]interface{}{
					"key":  "value",
					"key2": 23,
				},
			},
		},
	)
}

func TestComplexSwitchWorkflow(t *testing.T) {
	wf := testdata.GetWorkflowWithComplexSwitchTask()
	err := testdata.ValidateWorkflowRegistration(wf)
	if err != nil {
		t.Fatal(err)
	}
	receivedWf, _, err := testdata.MetadataClient.Get(context.Background(), wf.GetName(), nil)
	if err != nil {
		t.Fatal(err)
	}
	counter := countMultipleSwitchInnerTasks(receivedWf.Tasks...)
	assert.Equal(t, 7, counter)

	err = testdata.ValidateWorkflowDeletion(wf)
	if err != nil {
		t.Fatal(
			"Failed to delete workflow. Reason: ", err.Error(),
		)
	}
}

func countMultipleSwitchInnerTasks(tasks ...model.WorkflowTask) int {
	counter := 0
	for _, task := range tasks {
		counter += countSwitchInnerTasks(task)
	}
	return counter
}

func countSwitchInnerTasks(task model.WorkflowTask) int {
	counter := 1
	if task.Type_ != "SWITCH" {
		return counter
	}
	for _, value := range task.DecisionCases {
		counter += countMultipleSwitchInnerTasks(value...)
	}
	return counter
}
