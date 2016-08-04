import * as React from 'react';
import * as ReactDOM from 'react-dom';
import classNames from 'classnames';

export default class Todo extends React.Component {
    constructor(props) {
        super(props);

        this.ESCAPE_KEY = 27;
        this.ENTER_KEY = 13;
        this.state = {
            editText: this.props.todo.title
        };
    }

    handleSubmit() {
        const val = this.state.editText.trim();
        if (val) {
            this.props.onSave(val);
            this.setState({editText: val});
        } else {
            this.props.onDestroy();
        }
    }

    handleEdit() {
        this.props.onEdit();
        this.setState({editText: this.props.todo.title});
    }

    handleKeyDown(event) {
        if (event.which === this.ESCAPE_KEY) {
            this.setState({editText: this.props.todo.title});
            this.props.onCancel(event);
        } else if (event.which === this.ENTER_KEY) {
            this.handleSubmit();
        }
    }

    handleChange(event) {
        if (this.props.editing) {
            this.setState({editText: event.target.value});
        }
    }

    shouldComponentUpdate(nextProps, nextState) {
        return (
            nextProps.todo !== this.props.todo ||
            nextProps.editing !== this.props.editing ||
            nextState.editText !== this.state.editText
        );
    }

    componentDidUpdate(prevProps) {
        if (!prevProps.editing && this.props.editing) {
            const node = ReactDOM.findDOMNode(this.refs.editField);
            node.focus();
            node.setSelectionRange(node.value.length, node.value.length);
        }
    }

    render() {
        return (
            <li className={classNames({
              completed: this.props.todo.completed,
              editing: this.props.editing
            })}>
            <div className="view">
                <input
                    className="toggle"
                    type="checkbox"
                    checked={this.props.todo.completed}
                    onChange={this.props.onToggle}
                />
                <label onDoubleClick={this.handleEdit.bind(this)}>
                    {this.props.todo.title}
                </label>
                <button className="destroy" onClick={this.props.onDestroy} />
            </div>
            <input
                ref="editField"
                className="edit"
                value={this.state.editText}
                onBlur={this.handleSubmit.bind(this)}
                onChange={this.handleChange.bind(this)}
                onKeyDown={this.handleKeyDown.bind(this)}
            />
            </li>
        );
    }
}

Todo.propTypes = {
    editing: React.PropTypes.bool,
    editText: React.PropTypes.string,
    onToggle: React.PropTypes.func,
    onSave: React.PropTypes.func,
    onDestroy: React.PropTypes.func,
    onEdit: React.PropTypes.func,
    onCancel: React.PropTypes.func,
    todo: React.PropTypes.object
}
