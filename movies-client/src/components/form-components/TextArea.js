const TextArea = (props) => {
    return (
      <div className="mb-3">
        <label htmlFor="description" className="form-label">
          {props.title}
        </label>
        <textarea
          className="form-control"
          id={props.name}
          name={props.name}
          value={props.value}
          onChange={props.handleChange}
          rows={props.rows}
        />
      </div>
    );
  };
  
  export default TextArea;
  