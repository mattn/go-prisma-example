import {html, render} from 'https://unpkg.com/lit-html@3.1.3/lit-html.js';
import {classMap} from 'https://unpkg.com/lit-html/directives/class-map';
import ky from 'https://unpkg.com/ky@1.2.3/distribution/index.js';

const TaskItem = (
  task,
  onCheck,
  onClickRemove,
) => html`
  <div class=${classMap({ task: true, completed: task.completed })}>
    <input
      type="checkbox"
      @change=${(e) => onCheck(e.currentTarget.checked)}
      .checked=${task.completed}
    />
    <div class="text">${task.text}</div>
    <div class="remove" @click=${onClickRemove}>×</div>
  </div>
`

const TaskList = (
  tasks,
  onCheck,
  onClickRemove,
) => html`
  <div class="task-list">
    ${tasks.map(t =>
      TaskItem(t, checked => onCheck(t.id, checked), () => onClickRemove(t.id))
    )}
  </div>
`

const NewTask = (
  inputText,
  onInput,
  onKeyPress,
) => html`
  <div class="new-task">
    <input
      type="text"
      .value=${inputText}
      @input=${(e) => onInput(e.currentTarget.value)}
      @keypress=${onKeyPress}
      placeholder="what we have to do?"
    />
  </div>
`

const App = () => html`
  <div>
    <h1>ToDos</h1>
    ${NewTask(
      store().inputText,
      inputText => store({ inputText }),
      e => {
        if (e.key === "Enter") {
          (async () => {
            const state = store();
            const task = await ky.post('/tasks', {
              json: { text: state.inputText }
            }).json();
            store({
              tasks: [...state.tasks, task],
              inputText: ""
            })
          })();
        }
      }
    )}
    ${TaskList(
      store().tasks,
      (id, completed) => {
        (async () => {
          const task = await ky.post('/tasks/' + id, {
            json: { id: id, completed: completed }
          }).json();
          store({
            tasks: store().tasks.map(t => (t.id === id ? { ...t, completed } : t))
          })
        })()
      },
      id => {
        (async () => {
          const task = await ky.delete('/tasks/' + id, {
            json: { id: id }
          }).json();
          store({ tasks: store().tasks.filter(t => t.id !== id) })
        })()
      }
    )}
  </div>
`

const renderApp = () => render(App(), document.body)

const createStore = (initialState) => {
  let data = initialState

  return (update) => {
    if (update) {
      data = { ...data, ...update }
      renderApp()
    }
    return data
  }
}

let store = createStore({
  tasks: await ky.get('/tasks').json(),
  selectedTasks: [],
  inputText: "",
})

renderApp()
